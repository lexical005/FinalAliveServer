package main

import (
	"ffCommon/util"
	"fmt"
)

var mapSetupIAP = make(map[string]func(string, map[string]string) (string, error))

func onClientSetupIAP(reqClient string, dictData map[string]string) (string, error) {
	// 异常保护
	util.PanicProtect()

	if channel, ok := dictData["channel"]; !ok {
		return "", fmt.Errorf("onClientSetupIAP dictData not contain channel")
	} else if callback, ok := mapSetupIAP[channel]; !ok {
		return "", fmt.Errorf("onClientSetupIAP dictData invalid channel")
	} else {
		return callback(reqClient, dictData)
	}
}
