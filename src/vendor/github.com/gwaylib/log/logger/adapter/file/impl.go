package file

import (
	"fmt"
	"math"
	"os"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/proto"
)

type Adapter struct {
	fileName   string // 文件路径及名称
	backups    int32  // 保留的文件数
	fileSize   int64  // 每个文件大小
	file       *os.File
	buffer     chan *proto.Proto
	closeEvent chan bool
	doneEvent  chan bool
}

func (a *Adapter) Close() {
	a.closeEvent <- true
	<-a.doneEvent
	if a.file != nil {
		if err := a.file.Close(); err != nil {
			println(fmt.Sprintf("file adapter close failed: %s", err.Error()))
		}
	}
}

func (a *Adapter) Put(p *proto.Proto) {
	a.buffer <- p
	return
}

func (a *Adapter) run() {
	for {
		select {
		case p := <-a.buffer:
			// TODO: deal log
			for _, d := range p.Data {
				fmt.Println(d)
			}
		case <-a.closeEvent:
			println("closed")
			a.doneEvent <- true
			return
		}
	}
}

func New(fileName string) logger.Adapter {
	return NewFile(fileName, math.MaxInt32, 1024*1024)
}

func NewFile(fileName string, backups int32, fileSize int64) logger.Adapter {
	a := &Adapter{
		fileName:   fileName,
		backups:    backups,
		fileSize:   fileSize,
		buffer:     make(chan *proto.Proto, 100),
		closeEvent: make(chan bool, 1),
		doneEvent:  make(chan bool, 1),
	}
	go a.run()
	return a
}
