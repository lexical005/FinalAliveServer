package main

import (
	"ffCommon/log/log"
	"ffCommon/version"

	"net/http"
)

// 固定返回值
var hotfixDataInvalid = []byte("hotfix=3") // 版本无效, 必须与枚举值 HOTFIX_STATUS_INVALID 保持一致

// 响应手机端的 get 方式的 check 请求
// 提交参数:
//		model: 机型
//		channel: 渠道
//		version: 版本号
//		uuid: 全球唯一标识
// 返回值(可转换为 lua 表的字符串):
//	return {hotfix=0,list={"zip_file_1",size_1,"zip_file_2",size_2}}
//	return {hotfix=1}
//	return {hotfix=2,down_url="127.0.0.1:8088/wing.apk"}
//		hotfix: HOTFIX_STATUS
//		zip_file_1: 要下载的更新包
//		zip_size_1: 要下载的更新包的大小(字节)
//		down_url: 完整版本的下载地址

// 目前先简化处理, 必须提交渠道和版本号
//		channel: 渠道, 字符串
//		version: 版本号, 字符串1.1.101
func getCheck(rw http.ResponseWriter, req *http.Request) {
	// 异常保护
	defer func() {
		if r := recover(); r != nil {
			log.RunLogger.Println(r)
		}
	}()

	// 只响应 get 方式的请求
	if req.Method != "GET" {
		return
	}

	// 解析参数
	req.ParseForm()

	// 必须有 version, 且必须携带且只能携带一个在允许列表内的参数
	v1, ok1 := req.Form["version"]
	if !ok1 || len(v1) != 1 {
		// 未携带 version 参数, 或携带的 version 参数的值超过一个
		log.RunLogger.Println("not take version")
		return
	}

	userVersion, err := version.New(v1[0])
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	// 必须有 md5, 且必须携带且只能携带一个在允许列表内的参数
	var md5 string
	if v1, ok1 := req.Form["md5"]; !ok1 || len(v1) != 1 {
		// 未携带 md5 参数, 或携带的 md5 参数的值超过一个
	} else {
		md5 = v1[0]
	}

	// 必须有 channel, 且必须携带且只能携带一个指定的参数
	var response []byte

	v1, ok1 = req.Form["channel"]
	if !ok1 || len(v1) != 1 {
		// 未携带 channel 参数, 或携带的 channel 参数的值超过一个
		log.RunLogger.Println("not take channel")
		return
	}
	var channelName = v1[0]

	// 非开发渠道时，是不是已知的有效渠道
	if channelInfo, ok3 := globalChannelInfo[channelName]; !ok3 {
		log.RunLogger.Println("invalid channel:" + channelName)
		response = []byte("")
	} else {
		response = channelInfo.checkVersion(userVersion, md5)
	}

	rw.Write(response)

	log.RunLogger.Println(channelName, userVersion)
	log.RunLogger.Println(string(response))
	log.RunLogger.Println(req.RemoteAddr)
	log.RunLogger.Println("")
	log.RunLogger.Println("")
}
