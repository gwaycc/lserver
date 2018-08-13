
使用例子:

在具体项目中构建一个log包
```text
mkdir $PRJ_ROOT/src/log
cp log.go $PRJ_ROOT/src/log
```

自行修改以下配置信息

```text

// 日志配置信息
level   = proto.LevelDebug                                    // 日志输出的级别
ctx     = &proto.Context{"default", "1.0.0", logger.HostName} // 产生日志地方的客户端信息，用于日志服务器识别来源
adapter = []logger.Adapter{stdio.New(os.Stdout)}              // 日志输出适配器，用于输出日志

// 日志输出适配器，用于输出日志
// 日志服务器依赖于github.com/gwaycc/lserver项目
// adapter = []logger.Adapter{stdio.New(os.Stdout), lserver.New("127.0.0.1:11301", "log.gwaycc.com", 100)}
lg = New("default") // 系统默认输出前缀
```

配置自定log包后，调用例子:
例子1:
```text
import "log"

func main(){
    log.Info("Hello")
}
```

例子2:
``` text
import l "log"

var log = l.New("testing")

func main(){
    log.Info("Hello")
}
```

