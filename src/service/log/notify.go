package main

import (
	"fmt"
	"module/db/alarm"
	"os"
	"time"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/logger"
	gomail "gopkg.in/gomail.v2"
)

var (
	alertor = alarm.NewAlarm()
)

func init() {
	if _, err := alertor.LoadCfg(); err != nil {
		log.Info(errors.As(err))
	}
	cfg := alertor.Cfg()
	if err := alertor.Apply(cfg); err != nil {
		log.Info(errors.As(err))
	}
	go alertor.Deamon()
}

const (
	smsfmt = " md5: %s\n platform: %s\n level: %d"

	mailfmt = `
server-name : %+v
------------------------------------------------
md5         : %s
platform    : %s
version     : %s
ip          : %s
date        : %s
level       : %d
logger      : %s
msg         : %s`
)

// startTimes -- 起始触发的次数
func SendSms(data *DbTable, startTimes int) {
	platform := data.Platform()

	// 发送限制
	if !isSmsAble(platform, startTimes, time.Now(), 30*time.Minute) {
		return
	}

	msg := fmt.Sprintf(
		smsfmt,
		data.MD5()[:7],
		platform,
		data.Level(),
	)

	// TODO: send sms warnning
	_ = msg

}

// startTimes -- 起始触发的次数
func SendMail(data *DbTable, startTimes int) {
	platform := data.Platform()
	// 发送限制
	if !isMailAble(platform, startTimes, time.Now(), 30*time.Minute) {
		return
	}

	title := fmt.Sprintf(
		"log-server: %s",
		platform,
	)
	host, _ := os.Hostname()
	msg := fmt.Sprintf(
		mailfmt,
		host,
		data.MD5(),
		platform,
		data.Version(),
		data.Ip(),
		data.Date(),
		data.Level(),
		data.Logger(),
		data.MsgIndent(),
	)

	// get email receiver
	to := []string{}
	cfg := alertor.Cfg()
	for _, r := range cfg.Receivers {
		if len(r.Email) == 0 {
			continue
		}
		to = append(to, r.Email)
	}
	if len(to) > 0 {
		mailClient, ok := alertor.MailClient()
		if !ok {
			return
		}
		m := gomail.NewMessage()
		m.SetHeader("From", fmt.Sprintf("lserver-log<%s>", cfg.MailServer.AuthName))
		m.SetHeader("To", to...)
		m.SetHeader("Subject", title)
		m.SetBody("text/plain", msg)
		if err := mailClient.SendMail(m); err != nil {
			logger.FailLog(errors.As(err))
			return
		}
	}
	return
}

func SendLogFail(err error) {
	logger.FailLog(errors.As(err))
}
