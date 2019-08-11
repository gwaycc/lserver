package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"

	l "github.com/gwaylib/log"
)

var log = l.New("service/beh")

var closers = []io.Closer{}

func ListenExit(c io.Closer) {
	closers = append(closers, c)
}
func notifyExit() {
	for _, c := range closers {
		c.Close()
	}
}
func main() {

	fmt.Println("[ctrl+c to exit]")
	end := make(chan os.Signal, 2)
	signal.Notify(end, os.Interrupt, os.Kill)
	<-end

	notifyExit()
}
