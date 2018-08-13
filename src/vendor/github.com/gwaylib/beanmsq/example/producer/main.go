package main

import (
	"fmt"
	"time"

	"github.com/gwaylib/beanmsq"
)

func main() {
	addr := "127.0.0.1:11300"
	tube := "testing_tube"
	end := make(chan bool, 1)
	dealing := make(chan bool, 1)

	// 生产者
	p := beanmsq.NewProducer(100, addr, tube)
	eventSize := 50000000
	seed := time.Now().UnixNano()
	buffer := make(chan int64, 1000)
	for i := 1000; i > 0; i-- {
		go func() {
			for {
				in := <-buffer
				if err := p.Put([]byte(fmt.Sprintf("%d", seed+int64(in)))); err != nil {
					panic(err)
				}
				dealing <- true
			}
		}()
	}
	go func() {
		for i := eventSize; i > 0; i-- {
			buffer <- int64(i)
		}
	}()

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

	p.Close()
}
