package main

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	constUploadTypeCache    = "UploadCache"    // 上传日志
	constUploadTypeSetupIAP = "UploadSetupIAP" // 建立IAP订单
)

type funcClientRequestDispatcher func(string, string, map[string]string) error

var allClientRequestDispatcher = make(map[string]funcClientRequestDispatcher, 2)

func onClientRequest(rw http.ResponseWriter, req *http.Request) {
	// 只响应 POST 方式的请求
	if req.Method != "POST" {
		return
	}

	// 异常保护
	defer util.PanicProtect()

	// 成功
	result := "SUCCESS"
	defer func() {
		rw.Write([]byte(result))
	}()

	// 解析参数
	req.ParseForm()

	// Header携带的参数
	ReqClient := req.Header.Get("ReqClient")
	UploadType := req.Header.Get("UploadType")

	// 具体数据
	bodyContent, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	if UploadType == constUploadTypeCache {
		var listRequestDetailString []string
		if err := json.Unmarshal(bodyContent, &listRequestDetailString); err != nil {
			log.RunLogger.Println("unmarshal body failed", ReqClient, string(bodyContent), err)
			result = "ABANDON"
			return
		}

		for _, oneData := range listRequestDetailString {
			log.RunLogger.Println("upload data", ReqClient, oneData)

			dictData := make(map[string]string, 8)
			if err := json.Unmarshal([]byte(oneData), &dictData); err != nil {
				log.RunLogger.Println("unmarshal upload data failed", ReqClient, oneData, err)
				result = "ABANDON"
				continue
			}

			if ReqType, ok := dictData["ReqType"]; ok {
				if f, ok := allClientRequestDispatcher[ReqType]; ok {
					err := f(ReqClient, oneData, dictData)
					if err != nil {
						log.RunLogger.Println("response upload data failed", ReqClient, oneData, err)
						result = "FAIL"
					}
				} else {
					err := onClientRequestDefault(ReqClient, oneData, dictData)
					if err != nil {
						log.RunLogger.Println("default response request failed", ReqClient, oneData, err)
						result = "FAIL"
					}
				}
			} else {
				log.RunLogger.Println("upload data not contain ReqType", ReqClient, oneData)
				result = "FAIL"
			}
		}
	} else if UploadType == constUploadTypeSetupIAP {
		log.RunLogger.Println("upload data", ReqClient, string(bodyContent))

		var err error
		dictData := make(map[string]string, 8)
		if err = json.Unmarshal(bodyContent, &dictData); err != nil {
			log.RunLogger.Println("unmarshal upload data failed", ReqClient, string(bodyContent), err)
			result = "ABANDON"
			return
		}

		result, err = onClientSetupIAP(ReqClient, dictData)
		if err != nil {
			log.RunLogger.Println("setup iap failed", ReqClient, string(bodyContent), err)
			result = "FAIL"
			return
		}
	}
}
