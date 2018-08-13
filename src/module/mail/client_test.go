package mail

import (
	"testing"

	"gopkg.in/gomail.v2"
)

func TestOffice365Client(t *testing.T) {
	host := "smtp.163.com"
	userName := "free1139@163.com"
	userPwd := "" // 测试时请填写真实密码
	msg := gomail.NewMessage()
	msg.SetHeader("From", userName)
	msg.SetHeader("To", "free1139@163.com")
	msg.SetHeader("Subject", "Hello!")
	msg.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	mc := NewMailClient(host, 587, userName, userPwd)
	defer mc.Close()
	if err := mc.SendMail(msg); err != nil {
		t.Fatal(err)
	}
}
