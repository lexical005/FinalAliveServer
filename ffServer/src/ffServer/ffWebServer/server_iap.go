package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"io/ioutil"

	"net/http"
	"strings"
)

var serverIAPResponse = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<response>
    <packageid>
        <hret>0</hret>
        <message>Successful</message>
    </packageid>
</response>`)

func onServerRequestIAP(rw http.ResponseWriter, req *http.Request) {
	// 只响应 POST 方式的请求
	if req.Method != "POST" {
		return
	}

	// 异常保护
	defer util.PanicProtect()

	// 成功
	// result := "SUCCESS"
	defer func() {
		rw.Write([]byte(serverIAPResponse))
	}()

	// 解析参数
	req.ParseForm()

	content, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	s := string(content)
	s = strings.Replace(s, "\r\n", "", -1)
	s = strings.Replace(s, "\n", "", -1)

	log.RunLogger.Println("onServerRequestIAP:", s)
}
