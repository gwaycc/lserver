package logger

import (
	"os"
	"time"

	"github.com/gwaylib/log/logger/proto"
)

var (
	DefaultContext = proto.Context{"default", "0.0.0", HostName}
)

// logger for private
type Logger struct {
	context    *proto.Context
	adapters   []Adapter
	loggerName string
	level      proto.Level
}

func New(ctx *proto.Context, loggerName string, level proto.Level, adapter ...Adapter) *Logger {
	lg := &Logger{
		context:    ctx,
		adapters:   adapter,
		loggerName: loggerName,
		level:      level,
	}
	return lg
}

func (l *Logger) Debug(msg ...interface{}) {
	l.put(proto.LevelDebug, proto.ToMsg(msg...))
}

func (l *Logger) Debugf(f string, msg ...interface{}) {
	l.put(proto.LevelDebug, proto.ToMsgf(f, msg...))
}

func (l *Logger) Info(msg ...interface{}) {
	l.put(proto.LevelInfo, proto.ToMsg(msg...))
}

func (l *Logger) Infof(f string, msg ...interface{}) {
	l.put(proto.LevelInfo, proto.ToMsgf(f, msg...))
}

func (l *Logger) Warn(msg ...interface{}) {
	l.put(proto.LevelWarn, proto.ToMsg(msg...))
}

func (l *Logger) Warnf(f string, msg ...interface{}) {
	l.put(proto.LevelWarn, proto.ToMsgf(f, msg...))
}

func (l *Logger) Error(msg ...interface{}) {
	l.put(proto.LevelError, proto.ToMsg(msg...))
}

func (l *Logger) Errorf(f string, msg ...interface{}) {
	l.put(proto.LevelError, proto.ToMsgf(f, msg...))
}

func (l *Logger) Fatal(msg ...interface{}) {
	m := proto.ToMsg(msg...)
	l.put(proto.LevelFatal, m)
	l.Close()
	panic(m)
}

func (l *Logger) Fatalf(f string, msg ...interface{}) {
	m := proto.ToMsgf(f, msg...)
	l.put(proto.LevelFatal, m)
	l.Close()
	panic(m)
}

// Exit
// log an info level message, and close log, then call os.Exit(code)
//
// Param
// code -- code of os exit
// msg -- exit message
func (l *Logger) Exit(code int, msg ...interface{}) {
	m := proto.ToMsg(msg...)
	l.put(proto.LevelInfo, m)
	l.Close()
	os.Exit(code)

}

func (l *Logger) put(level proto.Level, msg []byte) {
	// ignore
	if level < l.level {
		return
	}
	l.Put([]*proto.Data{&proto.Data{time.Now(), level, l.loggerName, msg}})
}

func (l *Logger) Put(data []*proto.Data) {
	if len(data) == 0 {
		panic("len(data) == 0")
	}

	log := &proto.Proto{
		l.context,
		data,
	}

	// send log
	for _, adapter := range l.adapters {
		adapter.Put(log)
	}
}

// Close logger
// it is over time, suggest to call Exit because when logger closed, program also done.
func (l *Logger) Close() {
	for _, a := range l.adapters {
		a.Close()
	}
}
