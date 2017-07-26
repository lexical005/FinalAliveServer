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

// PanicProtect if panic then recover and print panic stack
func PanicProtect(extras ...interface{}) {
	if x := recover(); x != nil {
		PrintPanicStack(x, extras)
	}
}

// SafeGo 保护执行一个goroutine, 即使goroutine内部执行时发生异常, 也不会导致整个应用程序退出
// f: 要保护执行的函数
// params: 要传递的参数
func SafeGo(f func(params ...interface{}), params ...interface{}) {
	defer func() {
		if x := recover(); x != nil {
			PrintPanicStack(x, params)
		}
	}()

	f(params...)
}
