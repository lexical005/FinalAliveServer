package main

import (
	"ffCommon/log/log"

	"net/http"
	"strconv"
	"strings"
)

// 响应手机端的 get 方式的 remoteres 请求
// 提交参数:
//		remoteres: 要下载的远端文件
func getRemoteRes(rw http.ResponseWriter, req *http.Request) {
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

	// 必须有 remoteres, 且必须携带且只能携带一个参数
	var remoteResName string
	v, ok := req.Form["remoteres"]
	if !ok || len(v) != 1 {
		// 未携带 remoteres 参数, 或携带的 remoteres 参数的值超过一个
		log.RunLogger.Println("not take remoteres")
		return
	}
	remoteResName = v[0]

	remoteResName = strings.ToUpper(remoteResName)
	hotResInfo, ok := globalRemoteResCache[remoteResName]
	if !ok {
		genRemoteResMap()
		hotResInfo, ok = globalRemoteResCache[remoteResName]
		if !ok {
			log.RunLogger.Printf("remoteResName[%s] not found", remoteResName)
		}
	}

	if ok {
		rw.Header().Set("Content-Type", "application/octet-stream")
		rw.Header().Set("Accept-Ranges", "bytes")
		rw.Header().Set("Content-Length", strconv.Itoa(len(hotResInfo.buf)))

		rw.WriteHeader(http.StatusOK)

		rw.Write(hotResInfo.buf)

		log.RunLogger.Printf("remoteResName[%s] found, size[%d]", remoteResName, len(hotResInfo.buf))
	}

	log.RunLogger.Println(req.RemoteAddr)
	log.RunLogger.Println(req.RequestURI)
	log.RunLogger.Println("\n")
}

// 响应手机端的 get 方式的 remoteres 请求
// 提交参数:
//		res: 要下载的远端文件
func getRemoteResObsolete(rw http.ResponseWriter, req *http.Request) {
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

	// 必须有 res, 且必须携带且只能携带一个参数
	var remoteResName string
	v, ok := req.Form["res"]
	if !ok || len(v) != 1 {
		// 未携带 res 参数, 或携带的 res 参数的值超过一个
		log.RunLogger.Println("not take res")
		return
	}
	remoteResName = v[0]

	remoteResName = strings.ToUpper(remoteResName)
	hotResInfo, ok := globalRemoteResCache[remoteResName]
	if !ok {
		genRemoteResMap()
		hotResInfo, ok = globalRemoteResCache[remoteResName]
		if !ok {
			log.RunLogger.Printf("remoteResName[%s] not found", remoteResName)
		}
	}

	if ok {
		rw.Header().Set("Content-Type", "application/octet-stream")
		rw.Header().Set("Accept-Ranges", "bytes")
		rw.Header().Set("Content-Length", strconv.Itoa(len(hotResInfo.buf)))

		rw.WriteHeader(http.StatusOK)

		rw.Write(hotResInfo.buf)

		log.RunLogger.Printf("remoteResName[%s] found, size[%d]", remoteResName, len(hotResInfo.buf))
	}

	log.RunLogger.Println(req.RemoteAddr)
	log.RunLogger.Println(req.RequestURI)
	log.RunLogger.Println("\n")
}
