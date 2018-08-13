//
// 错误记录器
// 本设计补充并实现了系统的error接口，
// 用于发生错误时附带发生错误的原因、位置等信息以便还原现场。
// 因本错误设计含有比较大的数据量信息，因此需要注意被调用的频率，以避免影响到系统效率。
//
// 使用例子
// package main
//
// import "github.com/gwaylib/errors"
//
// func fn1(a int) error {
//    if a == 1{
//        return errors.ErrNoData.As(a)
//    }
//    return errors.New("not implements").As(a)
// }
//
// func fn2(b int) error {
//    return errors.As(fn1(b))
// }
//
// func main() {
//    err := fn2(2)
//    if err != nil {
//		  // 注意，errors.ErrNoData == err 会不一定成立, 需要使用Equal方法
//        if !errors.ErrNoData.Equal(err) {
//            panic(err)
//	      }
//        // 处理错误码相等的情况
//        fmt.Println(err)
//	  }
// }
//
package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

var (
	ErrNoData = New("data not found")
)

type Error interface {
	// 实现error接口
	Error() string
	// 实现json序列化接口
	MarshalJSON() ([]byte, error)

	// 返回New创建的码值
	Code() string

	// 记录当前调用的栈信息，并返回带新栈的新Error
	As(arg ...interface{}) Error

	// 将err转为本包的Error，再进行Code值比较
	// 等价于err1.Code() == err2.Code()
	Equal(err error) bool
}

//
// 比较两个错误的值是否相等.
// 该比较有两个范围，
// 一个是内存是否相等，常用于同个程序中产生的错误的比较，若内存相等，则两个错误是相等的;
// 一个是值是否相等，常用于跨程序中产生的错误比较，若不是此接口的Error，则比较Error()的值是否相等；
// 若属于此Error接口, 则比较Code()的值是否相等，
//
// 参数
// err1 -- 错误1
// err2 -- 另一个需要比较的错误
//
// 返回
// 返回是否相等，true相等
func Equal(err1 error, err2 error) bool {
	return equal(err1, err2)
}

func equal(err1 error, err2 error) bool {
	// Memory compare
	if err1 == err2 {
		return true
	}
	if err1 == nil || err2 == nil {
		return false
	}

	eImpl1, eImpl2 := ParseError(err1), ParseError(err2)
	return eImpl1.Code() == eImpl2.Code()
}

type ErrData struct {
	Err *string    `json:"Err"`
	As  [][]string `json:"As"`
}

type errImpl struct {
	data *ErrData
}

//
// 创建一个本包Error接口的实例
//
// 参数
// code -- 错误码或者文字描述，此值将用于Equal的比较
//
// 返回
// 返回一个新的Error实例
func New(code string) Error {
	return &errImpl{
		&ErrData{
			Err: &code,
			As:  [][]string{{caller(2), "[init]"}},
		},
	}
}

//
// 解析一个错误文本
// 通常它是一个本包的Error()序列化数据, 该数据是一个json数据，将直接被序列化为本包的Error类型;
// 若非本包的接口的序列化结构，将被直接New一个新的Error出来
//
// 参数
// src -- 需要解析的文本
//
// 返回
// 返回Error实例
func Parse(src string) Error {
	if len(src) == 0 {
		return nil
	}
	return parse(src)
}

//
// 将一个标准的错误转为本包的Error接口类型
// 若该错误本已经是本包的Error类型，则直接转为本包的Error并原样返回;
// 若该错误是非本包的Error类型，则调用error.Error()进行值解析创建一个新的本包Error
//
// 参数
// src -- 错误来源
//
// 返回
// 返回Error实例
func ParseError(src error) Error {
	if src == nil {
		return nil
	}
	if e, ok := src.(*errImpl); ok {
		return e
	}
	return parse(src.Error())
}

//
// 给一个错误构建错误定位信息
// 解析error时等价于ParseError，并在解析出的Error后构建当前置的错误定位信息。
// 若解析出的是本包的Error类型的实现，将在原实现基础上构建错误定位信息，此时也等价于Error的As方法调用。
//
// 参数
// err -- 任意类型的error实现
// reason -- 错误的原因，通常是引起发生错误的参数，以便记录并还原出发生错误时的调用。
//
// 返回
// 返回增加了定位信息的新Error实现, 因为是新实例返回，因此可以支持并发操作
func As(err error, reason ...interface{}) Error {
	if err == nil {
		return nil
	}
	e := ParseError(err).(*errImpl)
	as := []string{caller(2)}
	if len(reason) > 0 {
		as = append(as, fmt.Sprintf("%+v", reason))
	}
	return &errImpl{
		&ErrData{
			Err: e.data.Err,
			As:  append(e.data.As, as),
		},
	}
}

func parse(src string) *errImpl {
	if len(src) == 0 {
		return nil
	}
	if src[:1] != "{" {
		return New(src).(*errImpl)
	}

	data := &ErrData{}
	if err := json.Unmarshal([]byte(src), data); err != nil {
		return New(src).(*errImpl)
	}
	return &errImpl{data}
}

// call for domain
func caller(depth int) string {
	at := ""
	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		at = "domain of caller is unknown"
	}
	me := runtime.FuncForPC(pc)
	if me == nil {
		at = "domain of call is unnamed"
	}

	fileFields := strings.Split(file, "/")
	if len(fileFields) < 1 {
		at = "domain of file is unnamed"
		return at
	}
	funcFields := strings.Split(me.Name(), "/")
	if len(fileFields) < 1 {
		at = "domain of func is unnamed"
		return at
	}

	fileName := strings.Join(fileFields[len(fileFields)-1:], "/")
	funcName := strings.Join(funcFields[len(funcFields)-1:], "/")

	at = fmt.Sprintf("%s(%s:%d)", funcName, fileName, line)
	return at
}

// Error的Code方法实现
func (e *errImpl) Code() string {
	return *e.data.Err
}

// Error的Error方法实现
func (e *errImpl) Error() string {
	data, err := json.Marshal(e.data)
	if err != nil {
		return fmt.Sprintf("%+v", e.data)
	}
	return string(data)
}

// Error的MarshalJson方法实现
func (e *errImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.data)
}

// Error的Equal方法实现
func (e *errImpl) Equal(l error) bool {
	return equal(e, l)
}

// Error的As方法实现
func (e *errImpl) As(reason ...interface{}) Error {
	as := []string{caller(2)}
	if len(reason) > 0 {
		as = append(as, fmt.Sprintf("%+v", reason))
	}
	return &errImpl{
		&ErrData{
			Err: e.data.Err,
			As:  append(e.data.As, as),
		},
	}
}
