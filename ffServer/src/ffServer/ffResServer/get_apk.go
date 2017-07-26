package main

import (
	"ffCommon/log/log"

	"net/http"
	"strconv"
)

// 响应手机端的 get 方式的 apk 请求
// 提交参数:
//		zip_file: 要下载的更新包
func getApk(rw http.ResponseWriter, req *http.Request) {
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

	// 必须有 apk, 且必须携带且只能携带一个参数
	var subpackageName string
	v, ok := req.Form["apk"]
	if !ok || len(v) != 1 {
		// 未携带 apk 参数, 或携带的 apk 参数的值超过一个
		log.RunLogger.Println("not take apk")
		return
	}
	subpackageName = v[0]

	if apkSubpackage, ok := globalSubPackageCache[subpackageName]; ok {
		rw.Header().Set("Content-Type", "application/octet-stream")
		rw.Header().Set("Accept-Ranges", "bytes")
		rw.Header().Set("Content-Length", strconv.Itoa(apkSubpackage.fileSize))

		rw.WriteHeader(http.StatusOK)

		rw.Write(apkSubpackage.fileBuffer)
	}

	log.RunLogger.Println(subpackageName)
	log.RunLogger.Println(req.RemoteAddr)
	log.RunLogger.Println(req.RequestURI)
	log.RunLogger.Println("\n")
}
