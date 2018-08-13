# 使用例子

错误记录器

本设计补充并实现了系统的error接口，

用于发生错误时附带发生错误的原因、位置等信息以便还原现场。

因本错误设计含有比较大的数据量信息，因此需要注意被调用的频率，以避免影响到系统效率。

使用例子
```text
package main

import "github.com/gwaylib/errors"

func fn1(a int) error {
   if a == 1{
       return errors.ErrNoData.As(a)
   }
   return errors.New("not implements").As(a)
}

func fn2(b int) error {
   return errors.As(fn1(b))
}

func main() {
   err := fn2(2)
   if err != nil {
 	  // 注意，errors.ErrNoData == err 会不一定成立, 需要使用Equal方法
       if !errors.ErrNoData.Equal(err) {
           panic(err)
       }
       // 处理错误码相等的情况
       fmt.Println(err)
   }
}
```

输出用例解读
```text
输出：
{"Err":"test","As":[["errors.TestError(errors_test.go:113)","[init]"],["errors.TestError(errors_test.go:118)","[123 456]"]]}

Err -- errors.New()输入的Code值

As 
    -- ["errors.TestError(errors_test.go:113)","init"] 第一次被As调用记录的内容(初始化实例时第一次自动被调用)
        -- "errors.TestError(errors_test.go:113)" 第一次被As调用的位置信息
           -- errors.TestError 被调用了errors.As的函数 
           -- (errors_test.go:113) 被调用的文件行数
       -- "init" errors.As输入的所记录的第一个参数

    -- ["errors.TestError(errors_test.go:118)",123,456]] 第二次被As调用记录的内容
       -- "errors.TestError(errors_test.go:118)" 第二次被As调用的位置信息
            -- errors.TestError 被调用了errors.As的函数 
            -- (errors_test.go:118) 被调用的文件行数
       -- 123 errors.As输入的所记录的第一个参数
       -- 456 errors.As输入的所记录的第二个参数
       
```


# GO错误处理建议

## 优先处理错误，再处理正常逻辑, 错误将不容易被忽略。

```text
// 建议写法
rows, err := db.Query(...)
if err != nil{
    return errors.As(err)
}
defer rows.Close()
...

// 不建议写法
rows, err := db.Query(...)
if err == nil{
    defer rows.Close()
    // ...
} else {
    // 处理错误，或不处理错误
}
```


## 异常逻辑应尽快结束返回，结束函数流程，而不应再执行逻辑。

## 正常逻辑尽可能不放在if中，以使正常逻辑不缩进以便阅读。

``` text
// 除非必要，不建议使用此种错误类别的定义, 它易将错误定义发散到整个函数而失去可控性
func Get() (err error){
    defer func(){
        if err != nil{
            // 处理错误
        }
    }()

    // ...
}
```

## 除非明确了处理结果，否则错误应总是向上返回给调用者；
## 除非明确了处理结果，调用者不应丢弃任何错误, 应给出用户提示或记录日志，以便完整知道系统发生了什么。


