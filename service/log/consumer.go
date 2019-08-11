package main

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"lserver/module/etc"

	"github.com/gwaylib/beanmsq"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/logger/proto"
)

func init() {
	c := beanmsq.NewConsumer(etc.Etc.String("beanstalk", "addr"), etc.Etc.String("beanstalk", "tube-log"))
	go c.Reserve(20*time.Minute, Handle)
	ListenExit(c)
}

func Handle(ctx context.Context, job *beanmsq.Job, tried int) bool {
	p, err := proto.Unmarshal(job.Body)
	if err != nil {
		log.Error(errors.As(err))
		return true
	}

	for _, d := range p.Data {
		msg := ""
		if utf8.Valid(d.Msg) {
			msg = string(d.Msg)
		} else {
			msg = fmt.Sprintf("%#v", d.Msg)
		}
		tb := NewDbTable(
			p.Context.Platform,
			p.Context.Version,
			p.Context.Ip,
			d.Date.Format(time.RFC3339Nano),
			d.Level.Int(),
			d.Logger,
			msg,
		)
		if err := InsertLog(tb); err != nil {
			log.Error(errors.As(err))
			return true
		}

		// do notify
		switch d.Level {
		case proto.LevelFatal:
			SendMail(tb, 0)
			SendSms(tb, 0)
		case proto.LevelError:
			SendMail(tb, 0)
			SendSms(tb, 0)
		case proto.LevelWarn:
			SendMail(tb, 20)
			SendSms(tb, 20)
		}

	}

	return true
}
