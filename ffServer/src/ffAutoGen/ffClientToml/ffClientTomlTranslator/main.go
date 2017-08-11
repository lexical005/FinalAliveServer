package main

import (
	"ffCommon/util"
	"ffLogic/ffGrammar"
)

var allTrans = make([]func(), 0, 10)
var allRead = make([]func(), 0, 1)

func transGrammar(grammar ffGrammar.Grammar) *Grammar {
	return &Grammar{
		Grammar: grammar.Origin(),
	}
}

// 将toml数据转换为pb字节流
func trans() {
	for _, f := range allTrans {
		f()
	}
}

func reads() {
	for _, f := range allRead {
		f()
	}
}

func main() {
	defer util.PanicProtect()

	reads()

	trans()
}
