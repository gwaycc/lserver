package proto

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/gommon/color"
)

const (
	LevelDebug = Level(0)
	LevelInfo  = Level(1)
	LevelWarn  = Level(2)
	LevelError = Level(3)
	LevelFatal = Level(4)
)

type Level int

func (l Level) Int() int {
	return int(l)
}

func (l Level) ColorString() string {
	switch l {
	case LevelDebug:
		return color.White("DEBUG")
	case LevelInfo:
		return color.Green("INFO ")
	case LevelWarn:
		return color.Yellow("WARN ")
	case LevelError:
		return color.Red("ERROR")
	case LevelFatal:
		return color.Red("FATAL")
	}
	return color.Reset(strconv.Itoa(int(l)))
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	}
	return strconv.Itoa(int(l))
}

type Context struct {
	Platform string `json:"platform"`
	Version  string `json:"version"`
	Ip       string `json:"ip"`
}

func (c Context) String() string {
	return fmt.Sprintf("platform:%s,version:%s,ip:%s", c.Platform, c.Version, c.Ip)
}

type Data struct {
	Date   time.Time `json:"date"`
	Level  Level     `json:"level"`
	Logger string    `json:"logger"`
	Msg    []byte    `json:"msg"`
}

func (d Data) String() string {
	return fmt.Sprintf("%s,level:%d,logger:%s,msg:%s", d.Date.Format(time.RFC3339), d.Level, d.Logger, string(d.Msg))
}

type Proto struct {
	Context *Context `json:"context"`
	Data    []*Data  `json:"data"`
}

func Unmarshal(src []byte) (*Proto, error) {
	log := &Proto{}
	if err := json.Unmarshal(src, log); err != nil {
		return nil, err
	}
	return log, nil
}

func Marshal(l *Proto) ([]byte, error) {
	data, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 生成msg的方法
func ToMsg(i ...interface{}) []byte {
	return []byte(fmt.Sprintln(i...))
}
func ToMsgf(f string, i ...interface{}) []byte {
	return ToMsg(fmt.Sprintf(f, i...))
}
