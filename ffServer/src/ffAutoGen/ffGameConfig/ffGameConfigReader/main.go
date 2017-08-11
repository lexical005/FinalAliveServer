package main

import (
	"ffCommon/util"
)

var allRead = make([]func(), 0, 1)

func reads() {
	for _, f := range allRead {
		f()
	}
}

func main() {
	defer util.PanicProtect()

	reads()
}
