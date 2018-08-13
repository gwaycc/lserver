package proto

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gwaylib/errors"
)

var testDataString = []*Data{
	// 空值测试
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg:    ToMsg(nil),
	},

	// 字符串测试
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg:    ToMsg("1"),
	},
	// 整型测试
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg:    ToMsg(1),
	},
	// 浮点测试
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg:    ToMsg(1.0),
	},
	// 字节测试
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg:    ToMsg([]byte{0xff}),
	},
	// 编码测试
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg:    ToMsg(string([]byte{0xff})),
	},

	// 哈希表测试
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg: ToMsg(map[string]string{
			"ok": "ok",
		}),
	},
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg: ToMsg(map[int]string{
			1: "ok",
		}),
	},
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg: ToMsg(map[int]int{
			1: 1,
		}),
	},

	// 接口测试
	&Data{
		Date:   time.Now(),
		Level:  LevelDebug,
		Logger: "test",
		Msg:    ToMsg(errors.New("test").As("test")),
	},
}

func TestDataString(t *testing.T) {
	for index, val := range testDataString {
		fmt.Printf("%d,%s\n", index, val.String())
	}
}

func TestMarshal(t *testing.T) {
	l := &LogProto{Context{"platform", "1.0.0", "127.0.0.0"}, testDataString}
	data, err := Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	newL, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if l.Context.String() != newL.Context.String() {
		t.Fatal(l.Context, newL.Context)
	}
	if len(l.Data) != len(newL.Data) {
		t.Fatal(len(l.Data), len(newL.Data))
	}
	for _, data := range newL.Data {
		fmt.Println(string(data.Msg))
	}
}

func TestToMsg(t *testing.T) {
	fmt.Println(string(ToMsg(1)))
	fmt.Println(string(ToMsg(1, 2)))
	fmt.Println(string(ToMsg("1")))
	fmt.Println(string(ToMsg("1", "2")))

	v := map[string]string{"1": "1"}
	fmt.Println(string(ToMsg("1", v)))

	log.Println("1", v)
	log.Println("1")
	log.Println(v)
}
