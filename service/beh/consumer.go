package main

import (
	"context"
	"time"

	"lserver/module/etc"

	"github.com/gwaylib/beanmsq"
	"github.com/gwaylib/log/behavior"
	"github.com/gwaylib/log/logger"
)

func init() {
	c := beanmsq.NewConsumer(etc.Etc.String("beanstalk", "addr"), etc.Etc.String("beanstalk", "tube-beh"))
	go c.Reserve(20*time.Minute, Handle)
	ListenExit(c)
}

func Handle(ctx context.Context, job *beanmsq.Job, tried int) bool {
	event, err := behavior.Parse(job.Body)
	if err != nil {
		logger.FailLog(err)
		return true
	}
	if err := insertBehavior(event); err != nil {
		logger.FailLog(err)
		return true
	}
	return true
}
