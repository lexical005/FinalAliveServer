package main

import (
	"ffCommon/log/log"

	"net/http"
	"strconv"
	"strings"
)

// 响应手机端的 get 方式的 hotres 请求
// 提交参数:
//		channel: 要下载的热更新资源所属的渠道
//		hotres: 要下载的热更新资源
func getHotRes(rw http.ResponseWriter, req *http.Request) {
	// 异常保护
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	// 只响应 get 方式的请求
	if req.Method != "GET" {
		return
	}

	// 解析参数
	req.ParseForm()

	// 必须有 channel, 且必须携带且只能携带一个参数
	v1, ok := req.Form["channel"]
	if !ok || len(v1) != 1 {
		// 未携带 channel 参数, 或携带的 channel 参数的值超过一个
		log.RunLogger.Println("not take channel")
		return
	}
	var channelName = v1[0]
	channelInfo, ok3 := globalChannelInfo[channelName]
	if !ok3 {
		log.RunLogger.Println("invalid channel:" + channelName)
		return
	}

	// 必须有 hotres, 且必须携带且只能携带一个参数
	v2, ok := req.Form["hotres"]
	if !ok || len(v2) != 1 {
		// 未携带 hotres 参数, 或携带的 hotres 参数的值超过一个
		log.RunLogger.Println("not take hotres")
		return
	}
	var hotResName = v2[0]

	hotResName = strings.ToUpper(hotResName)
	hotResInfo, ok := channelInfo.hotResMap[hotResName]
	if !ok {
		genRemoteResMap()
		hotResInfo, ok = channelInfo.hotResMap[hotResName]
		if !ok {
			log.RunLogger.Printf("hotResName[%s] not found", hotResName)
		}
	}

	if ok {
		rw.Header().Set("Content-Type", "application/octet-stream")
		rw.Header().Set("Accept-Ranges", "bytes")
		rw.Header().Set("Content-Length", strconv.Itoa(len(hotResInfo.buf)))

		rw.WriteHeader(http.StatusOK)

		rw.Write(hotResInfo.buf)

		log.RunLogger.Printf("hotResName[%s] found, size[%d]", hotResName, len(hotResInfo.buf))
	}

	log.RunLogger.Println(req.RemoteAddr)
	log.RunLogger.Println(req.RequestURI)
	log.RunLogger.Println("\n")
}
