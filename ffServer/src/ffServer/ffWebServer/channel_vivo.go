package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"ffCommon/log/log"
	"ffCommon/util"
	"ffLogic/ffDef"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// vivo侧的订单信息
type vivoSetupOrderInfo struct {
	RespCode int    `json:"RespCode"` // 成功返回：200，非200时，respMsg会提示错误信息。
	RespMsg  string `json:"respMsg"`  // 对应响应码的响应信息

	SignMethod    string `json:"signMethod"`    // 对关键信息进行签名的算法名称
	Signature     string `json:"signature"`     // 对关键信息签名后得到的字符串1，用于商户验签
	VivoSignature string `json:"vivoSignature"` // 对关键信息签名后得到的字符串2，vivoSDK使用
	VivoOrder     string `json:"vivoOrder"`     // vivo订单号
	OrderAmount   string `json:"orderAmount"`   // 单位：元，币种：人民币，必须精确到小数点后两位，如：1.01，100.00
}

// vivo侧的支付信息
type vivoPayInfo struct {
	RespCode int    `json:"RespCode"` // 成功返回：200，非200时，respMsg会提示错误信息。
	RespMsg  string `json:"respMsg"`  // 对应响应码的响应信息

	SignMethod  string `json:"signMethod"`  // 对关键信息进行签名的算法名称
	Signature   string `json:"signature"`   // 对关键信息签名后得到的字符串1，用于商户验签
	StoreID     int    `json:"storeId"`     // 定长20位数字，由vivo分发的唯一识别码
	StoreOrder  int    `json:"storeOrder"`  // 商户自定义，最长 64 位字母、数字和下划线组成
	VivoOrder   string `json:"vivoOrder"`   // vivo订单号
	OrderAmount string `json:"orderAmount"` // 单位：元，币种：人民币，必须精确到小数点后两位，如：1.01，100.00
	PayType     int    `json:"channel"`     // 用户使用的支付渠道
	VivoFee     string `json:"channelFee"`  // 渠道扣除的费用
}

type sdkChannelVIVO struct {
	name         string
	payKey       string
	fmtSignature string

	vivoServerURL  string
	vivoAppKey     string
	vivoAppID      string
	vivoCpID       string
	vivoVersion    string
	vivoSignMethod string
	vivoNotifyURL  string
}

func (vivo *sdkChannelVIVO) onSetupIAPvivo(reqClient string, dictData map[string]string) (string, error) {
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

		oriSignature := fmt.Sprintf(vivo.fmtSignature, orderAmount, orderDesc, orderTime, orderTitle, storeOrder)
		oriSignatureByte := md5.Sum([]byte(oriSignature))
		signature := strings.ToLower(hex.EncodeToString(oriSignatureByte[:]))

		vivo.insertSQLOnLocalSetupOrder(reqClient, signature, storeOrder, orderAmount)

		body, err := vivo.doPostOrderToVIVO(url.Values{
			"version":    {vivo.vivoVersion},
			"signMethod": {vivo.vivoSignMethod},
			"signature":  {signature},

			"storeId":    {vivo.vivoCpID},
			"appId":      {vivo.vivoAppID},
			"storeOrder": {storeOrder},
			"notifyUrl":  {vivo.vivoNotifyURL},

			"orderAmount": {orderAmount},
			"orderDesc":   {orderDesc},
			"orderTime":   {orderTime},
			"orderTitle":  {orderTitle},
		})
		if err != nil {
			return "", err
		}

		response, err := vivo.onVIVOSetupOrder(reqClient, body)
		if err != nil {
			return "", err
		}
		log.RunLogger.Println("vivo.onSetupIAPvivo signature", signature, "response", response)

		return response, nil
	}
}

// 发送订单信息到vivo, 以获取vivo侧建立的订单信息
func (vivo *sdkChannelVIVO) doPostOrderToVIVO(data url.Values) ([]byte, error) {
	resp, err := http.PostForm(vivo.vivoServerURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// vivo侧建立订单结果
func (vivo *sdkChannelVIVO) onVIVOSetupOrder(reqClient string, vivoResponse []byte) (string, error) {
	log.RunLogger.Println("vivo.onVIVOSetupOrder", reqClient, string(vivoResponse))

	// 解析vivosdk的返回值
	orderInfo := &vivoSetupOrderInfo{}
	if err := json.Unmarshal(vivoResponse, orderInfo); err != nil {
		err = fmt.Errorf("onSetupIAPvivo Unmarshal response[%s] failed:%v", string(vivoResponse), err)
		return "", err
	}

	// 如果返回成功, 则更新数据库
	if orderInfo.RespCode == 200 {
		vivo.updateSQLOnVIVOSetupOrder(reqClient, orderInfo)
	}

	// 转换为发送到客户端的数据
	result := "success"
	reason := ""
	if orderInfo.RespCode != 200 {
		result = "failed"
		reason = orderInfo.RespMsg
	}
	dictResponse := map[string]string{
		// 通用字段
		"result": result,
		"reason": reason,

		// vivo特殊字段
		"vivoSignature": orderInfo.VivoSignature,
		"vivoOrder":     orderInfo.VivoOrder,
	}
	response, _ := json.Marshal(dictResponse)

	return string(response), nil
}

//收到vivo的支付成功通知
func (vivo *sdkChannelVIVO) onPay(rw http.ResponseWriter, req *http.Request) {
	// 只响应 POST 方式的请求
	if req.Method != "POST" {
		return
	}

	// 异常保护
	defer util.PanicProtect()

	// 成功
	defer func() {
		rw.WriteHeader(200)
	}()

	// 解析参数
	req.ParseForm()

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	log.RunLogger.Println("vivo.onPay", string(body))

	payInfo := &vivoPayInfo{}
	if err := json.Unmarshal(body, payInfo); err != nil {
		err = fmt.Errorf("onVIVOServerIAP Unmarshal body[%s] failed:%v", string(body), err)
		return
	}

	vivo.updateSQLOnRecvVIVOPay(payInfo)
}

// 建立本地订单
func (vivo *sdkChannelVIVO) insertSQLOnLocalSetupOrder(reqClient, signature, storeOrder, orderAmount string) {
	funInsertOrderCallback := func(result ffDef.IDBQueryResult) {
		err := result.SQLResult()
		if err != nil {
			log.RunLogger.Printf("vivo.insertSQLOnLocalSetupOrder %s excute get error[%v]", result.SQL(), err)
		} else {
			count, err := result.RowsAffected()
			if err != nil {
				log.RunLogger.Printf("vivo.insertSQLOnLocalSetupOrder %s RowsAffected get error[%v]", result.SQL(), err)
			} else if count == 0 {
				log.RunLogger.Printf("vivo.insertSQLOnLocalSetupOrder %s RowsAffected count zero", result.SQL())
			}
		}
	}
	mysql.query(0, 100, funInsertOrderCallback, reqClient, signature, storeOrder, orderAmount)
}

// vivo端订单信息
func (vivo *sdkChannelVIVO) updateSQLOnVIVOSetupOrder(reqClient string, orderInfo *vivoSetupOrderInfo) {
	funcUpdateCallback := func(result ffDef.IDBQueryResult) {
		err := result.SQLResult()
		if err != nil {
			log.RunLogger.Printf("vivo.updateSQLOnVIVOSetupOrder %s excute get error[%v]", result.SQL(), err)
		} else {
			count, err := result.RowsAffected()
			if err != nil {
				log.RunLogger.Printf("vivo.updateSQLOnVIVOSetupOrder %s RowsAffected get error[%v]", result.SQL(), err)
			} else if count == 0 {
				log.RunLogger.Printf("vivo.updateSQLOnVIVOSetupOrder %s RowsAffected count zero", result.SQL())
			}
		}
	}
	mysql.query(0, 101, funcUpdateCallback, orderInfo.VivoOrder, orderInfo.Signature, reqClient)
}

// vivo端支付成功信息
func (vivo *sdkChannelVIVO) updateSQLOnRecvVIVOPay(payInfo *vivoPayInfo) {
	funcUpdateCallback := func(result ffDef.IDBQueryResult) {
		err := result.SQLResult()
		if err != nil {
			log.RunLogger.Printf("vivo.updateSQLOnRecvVIVOPay %s excute get error[%v]", result.SQL(), err)
		} else {
			count, err := result.RowsAffected()
			if err != nil {
				log.RunLogger.Printf("vivo.updateSQLOnRecvVIVOPay %s RowsAffected get error[%v]", result.SQL(), err)
			} else if count == 0 {
				log.RunLogger.Printf("vivo.updateSQLOnRecvVIVOPay %s RowsAffected count zero", result.SQL())
			}
		}
	}
	payOrderTime := time.Now().Format("2006-01-02 15:04:05")
	mysql.query(0, 102, funcUpdateCallback, payInfo.PayType, payInfo.VivoFee, payOrderTime, payInfo.Signature, payInfo.VivoOrder)
}

var vivo = &sdkChannelVIVO{
	name:   "vivo",
	payKey: "vivoPay",

	vivoServerURL:  "https://pay.vivo.com.cn/vivoPay/getVivoOrderNum",
	vivoAppKey:     "3a8ebac872f4a06711ae36d5beb85c8b",
	vivoAppID:      "567fddbe1662c8547c2c7672aeb7ec5a",
	vivoCpID:       "20150708100703361904",
	vivoVersion:    "1.0.0",
	vivoSignMethod: "MD5",
}

func init() {
	// Signature计算举例：
	// Signature=to_lower_case(md5_hex(appId=XXX&notifyUrl=XXX&orderAmount=XXX&orderDesc=XXX&orderTime=XXX&orderTitle=XXX&storeId=XXX&storeOrder=XXX&version=XXX&to_lower_case(md5_hex(App-key))))

	var vivoAppKeyMD5Byte = md5.Sum([]byte(vivo.vivoAppKey))
	var vivoAppKeyMD5ForUse = strings.ToLower(hex.EncodeToString(vivoAppKeyMD5Byte[:]))

	vivo.vivoNotifyURL = "http://" + appConfig.Net.DomainName + ":" + appConfig.Net.ListenPort + "/" + vivo.payKey

	vivo.fmtSignature = fmt.Sprintf("appId=%s&notifyUrl=%s&orderAmount=%%s&orderDesc=%%s&orderTime=%%s&orderTitle=%%s&storeId=%s&storeOrder=%%s&version=%s&%s",
		vivo.vivoAppID, vivo.vivoNotifyURL, vivo.vivoCpID, vivo.vivoVersion, vivoAppKeyMD5ForUse)

	mapSetupIAP[vivo.name] = vivo.onSetupIAPvivo
}
