package etc

import (
	"encoding/json"
	"io/ioutil"
	"module/mail"
	"os"
	"sync"
	"time"

	"github.com/gwaylib/conf"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
)

type MailServer struct {
	SmtpHost string `json:"stmp_host"`
	SmtpPort int    `json:"stmp_port"`
	AuthName string `json:"auth_name"`
	AuthPwd  string `json:"auth_pwd"`
}

// TODO:implement
type SmsServer struct {
	AuthName string `json:"auth_name"`
	AuthPwd  string `json:"auth_pwd"`
}

type Receiver struct {
	NickName string `json:"nickname"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
}

func SearchReceiver(arr []*Receiver, nickName string) (int, bool) {
	for i, a := range arr {
		if a.NickName == nickName {
			return i, true
		}
	}
	return -1, false
}

func RemoveReceiver(arr []*Receiver, i int) []*Receiver {
	switch i {
	case 0:
		return arr[1:]
	case len(arr) - 1:
		return arr[:i]
	default:
		result := []*Receiver{}
		result = append(result, arr[:i]...)
		result = append(result, arr[i+1:]...)
		return result
	}
}

type AlarmCfg struct {
	Readme     string      `json:"readme"`
	MailServer MailServer  `json:"mail_server"`
	SmsServer  SmsServer   `json:"sms_server"`
	Receivers  []*Receiver `json:"receviers"`
}

type Alarm struct {
	ticker     *time.Ticker
	mutex      sync.Mutex
	mailClient *mail.MailClient
	cfg        *AlarmCfg
}

func NewAlarm() *Alarm {
	s := &Alarm{}
	return s
}
func (s *Alarm) Deamon() {
	s.ticker = time.NewTicker(5 * time.Minute)
	for {
		<-s.ticker.C
		if _, err := s.LoadCfg(); err != nil {
			log.Warn(errors.As(err))
			continue
		}
		if err := s.Apply(s.Cfg()); err != nil {
			log.Warn(errors.As(err))
			continue
		}
	}
}
func (s *Alarm) LoadCfg() (*AlarmCfg, error) {
	cfg := &AlarmCfg{}
	fileName := conf.RootDir() + "/etc/alarm.cfg"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.As(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.As(err)
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, errors.As(err)
		}
	}

	s.cfg = cfg
	return cfg, nil
}

func (s *Alarm) SaveCfg() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.cfg.Readme) == 0 {
		s.cfg.Readme = "用于配置服务器的主动告警机制，每5分钟会被读取一次"
	}
	data, err := json.MarshalIndent(s.cfg, "", "	")
	if err != nil {
		return errors.As(err)
	}
	file := conf.RootDir() + "/etc/alarm.cfg"
	if err := ioutil.WriteFile(file, data, os.ModePerm); err != nil {
		return errors.As(err)
	}
	return nil
}

func (s *Alarm) Cfg() *AlarmCfg {
	if s.cfg.Receivers == nil {
		s.cfg.Receivers = []*Receiver{}
	}
	return s.cfg
}

func (s *Alarm) MailClient() (*mail.MailClient, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.mailClient, s.mailClient != nil
}

func (s *Alarm) Apply(cfg *AlarmCfg) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.mailClient != nil {
		s.mailClient.Close()
		s.mailClient = nil
	}
	if len(cfg.MailServer.SmtpHost) > 0 {
		s.mailClient = mail.NewMailClient(
			cfg.MailServer.SmtpHost,
			cfg.MailServer.SmtpPort,
			cfg.MailServer.AuthName,
			cfg.MailServer.AuthPwd,
		)
		if err := s.mailClient.Test(); err != nil {
			return errors.As(err)
		}
	}

	// TODO: make sms client
	s.cfg = cfg
	return nil
}

// When Apply or Deamon is called, need to call Close
func (s *Alarm) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.ticker != nil {
		s.ticker.Stop()
	}

	if s.mailClient != nil {
		s.mailClient.Close()
		s.mailClient = nil
	}
	return nil
}
