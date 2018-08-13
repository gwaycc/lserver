package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gwaylib/beanmsq"
)

func main() {
	addr := "127.0.0.1:11300"
	tube := "testing_tube"
	end := make(chan bool, 1)
	dealing := make(chan bool, 1)

	// 消费者
	c := beanmsq.NewConsumer(addr, tube)
	result := sync.Map{}
	_ = result
	handle := func(ctx context.Context, job *beanmsq.Job, tried int) bool {
		// 检查并发消息的安全性, 若出现重复，说明读取是不安全的
		// 注意，测试中发现并发读取时存在重复接收到数据的问题
		oldID, ok := result.LoadOrStore(string(job.Body), job.ID)
		if ok {
			panic(fmt.Sprintf("repeated:%d,%d,%s", oldID, job.ID, string(job.Body)))
		}
		fmt.Printf("receive job, tried:%d, job:%+v\n", tried, string(job.Body))
		dealing <- true
		return true
	}
	go c.Reserve(20*time.Minute, handle)
	// 等待消费者就绪
	time.Sleep(1e9)

	//	// 生产者
	//	p := beanmsq.NewProducer(100, addr, tube)
	//	eventSize := 50000000
	//	seed := time.Now().UnixNano()
	//	buffer := make(chan int64, 1000)
	//	for i := 1000; i > 0; i-- {
	//		go func() {
	//			for {
	//				in := <-buffer
	//				if err := p.Put([]byte(fmt.Sprintf("%d", seed+int64(in)))); err != nil {
	//					panic(err)
	//				}
	//			}
	//		}()
	//	}
	//	go func() {
	//		for i := eventSize; i > 0; i-- {
	//			buffer <- int64(i)
	//		}
	//	}()

	// 等待结束事件
	// 若1秒钟读不到数据，认为已经没有数据可读
	go func() {
		for {
			select {
			case <-time.After(1e9):
				end <- true
				return
			case <-dealing:
				// ignore
			}
		}
	}()
	<-end

	c.Close()
	// p.Close()
}
