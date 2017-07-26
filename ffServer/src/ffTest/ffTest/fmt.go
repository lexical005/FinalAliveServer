package main

import (
	"bytes"
	_ "ffCommon/util"
	"fmt"
	"os/exec"
	_ "path"
)

var fmtStr = `package ffConfig


import (
    "fmt"
)

type Mall struct {
   VipPackage []*VipPackage
}

func (M *Mall) String() string {
    s := ""
    s += "VipPackage"
    for _, row := range M.VipPackage {
        s += fmt.Sprintf("%v\n", row)
    }

    return s
}

type VipPackage struct {
    InfoInt  int
    InfoStr  string
    InfoIntSingle  []int
    InfoStrSingle  []string
    InfoIntMulti  []int
    InfoStrMulti  []string
}

func (V *VipPackage) String() string {
    s := "["
    s += fmt.Sprintf("InfoInt:%v,", V.InfoInt)
    s += fmt.Sprintf("InfoStr:%v,", V.InfoStr)
    s += fmt.Sprintf("InfoIntSingle:%v,", V.InfoIntSingle)
    s += fmt.Sprintf("InfoStrSingle:%v,", V.InfoStrSingle)
    s += fmt.Sprintf("InfoIntMulti:%v,", V.InfoIntMulti)
    s += fmt.Sprintf("InfoStrMulti:%v,", V.InfoStrMulti)
    s += "]"
    return s
}


`

func testGofmt() {
	saveFilePath := "mall.go"
	cmdStr := fmt.Sprintf("gofmt -w %v", saveFilePath)
	in := bytes.NewBuffer(nil)
	cmd := exec.Command(cmdStr)
	cmd.Stdin = in
	if err := cmd.Run(); err != nil {
		fmt.Printf("excute cmdStr[%v] get error[%v]", cmdStr, err)
	}
}
