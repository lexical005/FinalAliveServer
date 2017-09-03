package main

import (
	"ffAutoGen/ffGameConfig"
	"ffCommon/util"
)

func main() {
	defer util.PanicProtect(nil)

	ffGameConfig.ReadAllToml()
}
