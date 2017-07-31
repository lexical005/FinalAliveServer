package util

import (
	"fmt"

	"ffCommon/log/log"

	"github.com/davecgh/go-spew/spew"
)

// PrintPanicStack print panic stack
func PrintPanicStack(x interface{}, extras ...interface{}) {
	s := fmt.Sprintf("PrintPanicStack reason[%v]\n", x)
	for k := range extras {
		s += fmt.Sprintf("EXRA_INFO#%v DATA:%v\n", k, spew.Sdump(extras[k]))
	}
	log.FatalLogger.Println(s)
}

// PanicProtect if panic then recover and print panic stack. defer PanicProtect()
func PanicProtect(extras ...interface{}) {
	if x := recover(); x != nil {
		PrintPanicStack(x, extras...)
	}
}

// SafeGo 保护执行一个goroutine, 即使goroutine内部执行时发生异常, 也不会导致整个应用程序退出
// 	f: 保护执行的函数
//	c: 保护执行的函数执行完毕后, 安全清理函数, 允许为nil
//	params: 要传递的给f的参数
func SafeGo(f func(params ...interface{}), c func(), paramsF ...interface{}) {
	// 保护函数
	defer PanicProtect(paramsF...)

	// 保护执行的函数执行完毕后, 进行清理
	if c != nil {
		defer c()
	}

	// 保护执行的函数
	f(paramsF...)
}
