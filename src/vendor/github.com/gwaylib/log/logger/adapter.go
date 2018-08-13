package logger

import (
	"github.com/gwaylib/log/logger/proto"
)

type Adapter interface {
	Put(log *proto.Proto)
	Close()
}

// 处理适配器记录日志时产的错误，以println的方式输出至控制台。
// TODO: 考虑将错误写到系统日志当中。
func FailLog(err error) {
	println(err.Error())
}
