//
// Example
//
// import gopkg.in/gomail.v2
//
// m := gomail.NewMessage()
// m.SetHeader("From", "alex@example.com")
// m.SetHeader("To", "bob@example.com", "cora@example.com")
// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
// m.SetHeader("Subject", "Hello!")
// m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
// m.Attach("/home/Alex/lolcat.jpg")
//
// d := NewMailClient("smtp.example.com", 587, "user", "123456")
// defer d.Close()
//
// // Send the email to Bob, Cora and Dan.
// if err := d.SendMail(m); err != nil {
//     panic(err)
// }
//
package mail

import (
	"sync"

	"github.com/gwaylib/errors"
	gomail "gopkg.in/gomail.v2"
)

type MailClient struct {
	mux    sync.Mutex
	dialer *gomail.Dialer
	client gomail.SendCloser
}

func NewMailClient(host string, port int, username, passwd string) *MailClient {
	d := gomail.NewPlainDialer(host, port, username, passwd)
	return &MailClient{dialer: d}
}

func New(d *gomail.Dialer) *MailClient {
	return &MailClient{dialer: d}
}

func (mc *MailClient) SendMail(msg *gomail.Message) error {
	mc.mux.Lock()
	defer mc.mux.Unlock()
	if err := mc.sendMail(msg); err != nil {
		mc.close()
		return errors.As(err)
	}
	return nil
}

func (mc *MailClient) Test() error {
	mc.mux.Lock()
	defer mc.mux.Unlock()
	if err := mc.dial(); err != nil {
		return errors.As(err)
	}
	return nil
}

func (mc *MailClient) close() (err error) {
	if mc.client != nil {
		err = mc.client.Close()
		mc.client = nil
	}
	return errors.As(err)
}

func (mc *MailClient) dial() error {
	if mc.client != nil {
		return nil
	}
	c, err := mc.dialer.Dial()
	if err != nil {
		return errors.As(err)
	}

	mc.client = c
	return nil
}

func (mc *MailClient) sendMail(msg *gomail.Message) error {
	if err := mc.dial(); err != nil {
		return errors.As(err)
	}
	return errors.As(gomail.Send(mc.client, msg))
}

func (mc *MailClient) Close() error {
	mc.mux.Lock()
	defer mc.mux.Unlock()
	return mc.close()
}
