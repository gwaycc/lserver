package log

import (
	"os"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/adapter/stdio"
	"github.com/gwaylib/log/logger/proto"
)

var (
	// 日志配置信息
	level   = proto.LevelDebug                       // 日志输出的级别
	ctx     = &logger.DefaultContext                 // 产生日志地方的客户端信息，用于日志服务器识别来源
	adapter = []logger.Adapter{stdio.New(os.Stdout)} // 日志输出适配器，用于输出日志
	Log     = New("default")                         // 系统默认输出前缀
)

// 快速构建一个带前缀的日志
func New(prefix string) proto.Logger {
	return logger.New(ctx, prefix, level, adapter...)
}

// 值 0
// 调试信息,不提交服务器.
// 由平台控制,可控制台输出,相当于fmt.Print输出.
func Debug(msg ...interface{}) {
	Log.Debug(msg...)
}
func Debugf(f string, msg ...interface{}) {
	Log.Debugf(f, msg...)
}

// 等同于Info
func Print(v ...interface{}) {
	Info(v...)
}

// 等同于Info
func Printf(format string, v ...interface{}) {
	Infof(format, v...)
}

// 等同于Info
func Println(v ...interface{}) {
	Info(v...)
}

// 值 1
// 程序运行状态信息,提交服务器.
// 如启动、停止、重连等信息，体现了程序环境的变更状态。
func Info(msg ...interface{}) {
	Log.Info(msg...)
}
func Infof(f string, msg ...interface{}) {
	Log.Infof(f, msg...)
}

// 值 2
// 程序异常信息，提交服务器. 本类别不影响程序继续使用,但其结果可能会引出潜在的重大问题.
// 例如：请求的数据格式错误；网络连接错误(但重新连接后可自动修复), 连接超时等行为。
// 此类异常在一段时间如果出现过多，那么应该分析其中的原因，例如可能存在的问题：
// 被攻击、硬件老化、硬件达到了承载上限、对方服务出现异常等问题。
// 日志系统将发送一封邮件到相关人员。
func Warn(msg ...interface{}) {
	Log.Warn(msg...)
}
func Warnf(f string, msg ...interface{}) {
	Log.Warnf(f, msg...)
}

// 值 3
// 程序致命的错误信息， 提交服务器。此错误将影响到正常逻辑, 甚至平台因此而恐慌、停止服务的行为.
// 例如：数据库不可用、充值不可用、短信不可用、vos不可用等需要及时处理的行为都可定义为此类别。
// 日志系统将发送一封邮件、短信(或者其他实时联系方式)至相关人员
func Error(msg ...interface{}) {
	Log.Error(msg...)
}
func Errorf(f string, msg ...interface{}) {
	Log.Errorf(f, msg...)
}

// 值 4
// 检测到程序非正常结束。
// 日志系统将调用所有实时联系方式联系相关人员处理。
func Fatal(msg ...interface{}) {
	Log.Fatal(msg...)
}
func Fatalf(f string, msg ...interface{}) {
	Log.Fatalf(f, msg...)
}

// 退出操作，执行关闭日志操作并调用os.Exit
// 对于可执行程序来说，日志的退出意味着程序的退出。
//
// 参数
//  code -- 退出码值，由os.Exit调用
//  msg -- 记录的消息，级别是一个Info级别.
//
func Exit(code int, msg ...interface{}) {
	Log.Exit(code, msg...)
}
