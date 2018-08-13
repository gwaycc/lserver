//
// 例子
//
// c := NewConsumer("localhost:11130", "test")
//
// handle := func(ctx context.Context, job *beans.Job, tried int) bool{
//		// 处理结束后返回true删除数据
//		return true
// }
//
// 开启两个队列去并发读取
// go c.Reserve(20*time.Minute, handle)
// // go c.Reserve(20*time.Minute, handle) // 警告：目前beanstalkd消费者并发连接读取测试未通过, 多线程读取时注意去重操作
//
// // 在适当的地方关闭连接
// // c.Close() // Stop func to stop working
//
package beanmsq

import (
	"context"
	"io"
	"net"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/adapter/stdio"
	"github.com/gwaylib/log/logger/proto"
	beans "github.com/kr/beanstalk"
)

func IsErrNotFound(err error) bool {
	return strings.Index(strings.ToLower(err.Error()), "not found") >= 0
}

// 最大推送次数
const MAX_TRY_TIMES = 48 + 30 + 1

type Job struct {
	ID   uint64
	Body []byte
}

//
// 若发送不成功
// 返回true删除beanstalkd队件数据，否则不删除在一定时间后放回到就绪队中再次读取以便达到重试的效果。
// 重试机制分别间隔以1次3秒钟、30次每分钟、48次每小时再次尝试发送, 若48小时后未能发送成功，数据将被强制删除。
// 已放回就绪队列的次数通过tried进行了推送
type HandleContext func(ctx context.Context, job *Job, tried int) bool

type Consumer interface {
	io.Closer

	// reserveOut -- 与服务器连接的存活检查时间
	// handle -- 接收处理函数
	Reserve(reserveOut time.Duration, handle HandleContext) error
}

type consumer struct {
	addr     string
	tube     string
	workerMu sync.RWMutex
	isClosed bool
	workers  []io.Closer
}

func NewConsumer(addr, tube string) Consumer {
	c := &consumer{
		addr:    addr,
		tube:    tube,
		workers: []io.Closer{},
	}
	return c
}

func (c *consumer) Reserve(reserveOut time.Duration, handle HandleContext) error {
	c.workerMu.Lock()
	if c.isClosed {
		c.workerMu.Unlock()
		return errors.New("Consumer has closed")
	}
	w := newConsumer(c.addr, c.tube, handle, reserveOut)
	c.workers = append(c.workers, w)
	c.workerMu.Unlock()

	// 启动存活检查机制
	go w.aliveCheck()

	w.reserve()

	return nil
}
func (c *consumer) Close() error {
	c.workerMu.Lock()
	defer c.workerMu.Unlock()
	for _, w := range c.workers {
		w.Close()
	}
	return nil
}

type worker struct {
	// 日志器
	log proto.Logger
	// lock for Reader operator.
	mutex sync.Mutex

	// server addr
	addr string

	// channle name
	tubename string

	// handle which is pushed
	handle HandleContext

	// server connection
	conn *beans.TubeSet

	// work timeout for dealock
	reserveOut time.Duration

	// log error times
	connErrTimes int

	tryHistory map[uint64]int

	ctx  context.Context
	ctxc context.CancelFunc

	// signal command.
	sig_exit_reserve chan bool
	sig_alive        chan bool
	sig_end          chan bool
}

func newConsumer(addr, tube string, handle HandleContext, reserveOut time.Duration) *worker {
	return &worker{
		log:              logger.New(&logger.DefaultContext, tube, 0, stdio.New(os.Stderr)),
		mutex:            sync.Mutex{},
		addr:             addr,
		tubename:         tube,
		handle:           handle,
		reserveOut:       reserveOut,
		tryHistory:       make(map[uint64]int),
		sig_exit_reserve: make(chan bool, 1),
		sig_alive:        make(chan bool, 10),
		sig_end:          make(chan bool, 1),
	}
}

func (c *worker) timeout(msg string) {
	// 超时时关闭连接以便执行自动重连
	c.mutex.Lock()
	if c.ctxc != nil {
		c.ctxc()
	}
	c.disconn(msg)
	time.Sleep(1e9)
	c.mutex.Unlock()
}

func (c *worker) aliveCheck() {
	t := time.NewTimer(c.reserveOut)
	for {
		if !t.Stop() {
			t.Reset(c.reserveOut)
		} else {
			t.Stop()
			t = time.NewTimer(c.reserveOut)
		}

		select {
		case <-c.sig_alive:
			continue
		case <-t.C:
			c.timeout("alive check")
			continue
		}
	}
}

func (c *worker) reserve() {
	for {
		select {
		case <-c.sig_exit_reserve:
			c.sig_end <- true
			return
		default:
			c.mutex.Lock()
			if c.conn == nil {
				// 检查连接
				if err := c.connect(); err != nil {
					c.log.Warn(errors.As(err))
					c.mutex.Unlock()
					continue
				}
			}
			c.mutex.Unlock()

			// send alive signal
			c.sig_alive <- true

			// no lock
			id, body, err := c.conn.Reserve(c.reserveOut)
			if err != nil {
				if strings.Index(err.Error(), "use of closed network connection") >= 0 {
					// user manu close connection
					continue
				}

				c.log.Warn(errors.As(err))
				// 重连
				c.mutex.Lock()
				c.disconn(errors.As(err).Error())
				c.mutex.Unlock()
			}

			c.mutex.Lock()
			var ctx context.Context
			if c.ctxc != nil {
				c.ctxc()
			}
			ctx, c.ctxc = context.WithCancel(context.Background())

			// no lock
			c.do(ctx, &Job{id, body})
			c.mutex.Unlock()
		}
	}
}

func (c *worker) Close() error {
	c.sig_exit_reserve <- true

	c.mutex.Lock()
	c.disconn("normal shutdown")
	c.mutex.Unlock()

	<-c.sig_end
	return nil
}

func (c *worker) connect() error {
	// do close at first
	if c.conn != nil {
		c.disconn("to reconnect")
	}

	// connect
	c.log.Info("msq-c connect:" + c.tubename)
	kon, err := net.DialTimeout("tcp", c.addr, 20*1e9)
	if err != nil {
		c.connErrTimes++
		c.dealConnErrTimes(c.connErrTimes, errors.As(err))
		return errors.As(err)
	}
	c.conn = beans.NewTubeSet(beans.NewConn(kon), c.tubename)

	c.connErrTimes = 0
	return nil
}

func (c *worker) disconn(msg string) {
	if c.conn != nil {
		c.conn.Conn.Close()
		c.conn = nil
		c.log.Infof("msq-c closed:%s, by:%s", c.tubename, msg)
	}
}

// deal connction error
func (c *worker) dealConnErrTimes(times int, err error) {
	if err == nil {
		err = errors.New("unknow error")
	}
	switch times {
	case 10:
		// if error times equal 10, make a error log.
		c.log.Error(errors.As(err))
	case 3:
		// if error times equal 3, make warnning log.
		c.log.Warn(errors.As(err))
	}

	// if error times more than 10 times, do a sleep 30 sec.
	if times > 10 {
		time.Sleep(10 * 1e9)
	} else {
		time.Sleep(1 * 1e9)
	}
}

// do job
func (c *worker) do(ctx context.Context, job *Job) {
	times, _ := c.tryHistory[job.ID]
	defer func() {
		// recover for handle
		if r := recover(); r != nil {
			c.log.Error(errors.New("panic").As(r))
			debug.PrintStack()
			time.Sleep(10e9)
		}
	}()

	deal := c.handle(ctx, job, times)
	if deal {
		c.delJob(job)
	} else {
		c.nextTry(job)
	}
}

func (c *worker) nextTry(job *Job) {
	id := job.ID
	times, _ := c.tryHistory[id]
	times += 1
	c.tryHistory[id] = times

	// 若发送不成功
	// 分别间隔以1次3秒钟、30次每分钟、48次每小时再次尝试发送, 若48小时后未能发送成功，数据将被删除
	sleep := 10
	if times < 2 {
		sleep = 3 // 3 sec.
	} else if times < 30 {
		sleep = 60 // 1 minute
	} else if times < 78 {
		sleep = 60 * 60 // 1 hour
	} else {
		c.log.Warn(errors.New("delete data").As(string(job.Body)))
		// 48+6次后删除数据
		// delete job
		c.delJob(job)
		return
	}

	if err := c.conn.Conn.Release(job.ID, 0, time.Duration(sleep*1e9)); err != nil {
		if !IsErrNotFound(err) {
			c.log.Error(errors.As(err, job))
			c.disconn(err.Error())
		}
	}
	return
}

func (c *worker) delJob(job *Job) {
	id := job.ID
	delete(c.tryHistory, id)
	if err := c.conn.Conn.Delete(id); err != nil {
		if !IsErrNotFound(err) {
			log.Error(errors.As(err))
			c.disconn(err.Error())
		}
	}
}
