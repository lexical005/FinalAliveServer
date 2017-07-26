package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"ffCommon/log/log"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	vivoServerURL  = "https://pay.vivo.com.cn/vivoPay/getVivoOrderNum"
	vivoAppKey     = "3a8ebac872f4a06711ae36d5beb85c8b"
	vivoAppID      = "567fddbe1662c8547c2c7672aeb7ec5a"
	vivoCpID       = "20150708100703361904"
	vivoVersion    = "1.0.0"
	vivoSignMethod = "MD5"
	vivoNotifyURL  = "http://standalone.coola.tv:12580/vivoIAP"
)

var vivoKeys = []string{
	"appId",
	"notifyUrl",
	"orderAmount",
	"orderDesc",
	"orderTime",
	"orderTitle",
	"storeId",
	"storeOrder",
	"version",
}

type vivoSetupResponse struct {
	RespCode int    `json:"RespCode"` // 成功返回：200，非200时，respMsg会提示错误信息。
	RespMsg  string `json:"respMsg"`  // 对应响应码的响应信息

	SignMethod    string `json:"signMethod"`    // 对关键信息进行签名的算法名称
	Signature     string `json:"signature"`     // 对关键信息签名后得到的字符串1，用于商户验签
	VivoSignature string `json:"vivoSignature"` // 对关键信息签名后得到的字符串2，vivoSDK使用
	VivoOrder     string `json:"vivoOrder"`     // vivo订单号
	OrderAmount   string `json:"orderAmount"`   // 单位：元，币种：人民币，必须精确到小数点后两位，如：1.01，100.00
}

var vivoAppKeyMD5Byte = md5.Sum([]byte(vivoAppKey))
var vivoAppKeyMD5ForUse = strings.ToLower(hex.EncodeToString(vivoAppKeyMD5Byte[:]))

var vivoSignatureFmt = fmt.Sprintf("appId=%s&notifyUrl=%s&orderAmount=%%s&orderDesc=%%s&orderTime=%%s&orderTitle=%%s&storeId=%s&storeOrder=%%s&version=%s&%s",
	vivoAppID, vivoNotifyURL, vivoCpID, vivoVersion, vivoAppKeyMD5ForUse)

// Signature计算举例：
// Signature=to_lower_case(md5_hex(appId=XXX&notifyUrl=XXX&orderAmount=XXX&orderDesc=XXX&orderTime=XXX&orderTitle=XXX&storeId=XXX&storeOrder=XXX&version=XXX&to_lower_case(md5_hex(App-key))))

func onSetupIAPvivo(reqClient string, dictData map[string]string) (string, error) {
	orderTime := time.Now().Format("20060102150405")
	if orderAmount, ok := dictData["orderAmount"]; !ok || orderAmount == "" {
		return "", fmt.Errorf("onSetupIAPvivo not contain orderAmount")
	} else if orderDesc, ok := dictData["orderDesc"]; !ok {
		return "", fmt.Errorf("onSetupIAPvivo not contain orderDesc")
	} else if orderTitle, ok := dictData["orderTitle"]; !ok {
		return "", fmt.Errorf("onSetupIAPvivo not contain orderTitle")
	} else if storeOrder, ok := dictData["storeOrder"]; !ok || storeOrder == "" {
		return "", fmt.Errorf("onSetupIAPvivo not contain storeOrder")
	} else {
		if orderDesc == "" {
			orderDesc = "-"
		}
		if orderTitle == "" {
			orderTitle = "-"
		}

		oriSignature := fmt.Sprintf(vivoSignatureFmt, orderAmount, orderDesc, orderTime, orderTitle, storeOrder)
		oriSignatureByte := md5.Sum([]byte(oriSignature))
		signature := strings.ToLower(hex.EncodeToString(oriSignatureByte[:]))

		resp, err := http.PostForm(vivoServerURL, url.Values{
			"version":    {vivoVersion},
			"signMethod": {vivoSignMethod},
			"signature":  {signature},

			"storeId":    {vivoCpID},
			"appId":      {vivoAppID},
			"storeOrder": {storeOrder},
			"notifyUrl":  {vivoNotifyURL},

			"orderAmount": {orderAmount},
			"orderDesc":   {orderDesc},
			"orderTime":   {orderTime},
			"orderTitle":  {orderTitle},
		})
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// 解析vivosdk的返回值
		vivoSetupResponse := &vivoSetupResponse{}
		if err := json.Unmarshal(body, vivoSetupResponse); err != nil {
			err = fmt.Errorf("onSetupIAPvivo Unmarshal response[%s] failed:%v", string(body), err)
			return "", err
		}

		// 转换为发送到客户端的数据
		result := "success"
		reason := ""
		if vivoSetupResponse.RespCode != 200 {
			result = "failed"
			reason = vivoSetupResponse.RespMsg
		}
		dictResponse := map[string]string{
			// 通用字段
			"result": result,
			"reason": reason,

			// vivo特殊字段
			"vivoSignature": vivoSetupResponse.VivoSignature,
			"vivoOrder":     vivoSetupResponse.VivoOrder,
		}
		response, _ := json.Marshal(dictResponse)

		s := string(response)
		log.RunLogger.Println("onSetupIAPvivo signature", signature, "response", s)
		return s, nil
	}
}

func testVIVO() {
	dictDatas := map[string]string{
		"orderAmount": "0.01",
		"orderDesc":   "orderDesc",
		"orderTitle":  "orderTitle",
		"storeOrder":  "storeOrder",
	}
	fmt.Println(onSetupIAPvivo("test", dictDatas))
}
