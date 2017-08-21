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

// PanicProtect if panic then recover and print panic stack.
//	usage: defer PanicProtect()
func PanicProtect(onProtectEnd func(isPanic bool), extras ...interface{}) {
	isPanic := false
	if x := recover(); x != nil {
		PrintPanicStack(x, extras...)
		isPanic = true
	}

	if onProtectEnd != nil {
		onProtectEnd(isPanic)
	}
}

// SafeGo 保护执行一个goroutine, 即使goroutine内部执行时发生异常, 也不会导致整个应用程序退出
// 	f: 保护执行的函数
//	c: 保护执行的函数执行完毕后, 安全清理函数, 允许为nil
//	params: 要传递的给f的参数
func SafeGo(f func(params ...interface{}), c func(isPanic bool), paramsF ...interface{}) {
	isPanic := false

	// 保护执行的函数执行完毕后, 进行清理
	if c != nil {
		// 执行安全清理函数
		defer func() {
			// 保护函数
			defer PanicProtect(nil)

			// 安全清理函数
			c(isPanic)
		}()
	}

	// 保护执行的函数
	isPanic = doF(f, paramsF...)
}

func doF(f func(params ...interface{}), paramsF ...interface{}) (isPanic bool) {
	defer func() {
		if x := recover(); x != nil {
			PrintPanicStack(x, paramsF...)
			isPanic = true
		}
	}()

	// 保护执行的函数
	f(paramsF...)
	return
}
