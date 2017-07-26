package encrypt

import (
	"encoding/base64"
	"strings"
)

// EncodeToBase64 将字节流转换为可视化的base64编码返回，标准方式
// datas: 待编码的字节流
// result: 1 md5值转换为大写；-1 md5值转换为小写；0 md5值不做整理
func EncodeToBase64(datas []byte, result int) string {
	base64Result := base64.StdEncoding.EncodeToString(datas)
	if result == 1 {
		return strings.ToUpper(base64Result)
	} else if result == -1 {
		return strings.ToLower(base64Result)
	}
	return base64Result
}

// DecodeFromBase64 将base64编码后的字符串解码为原始字节流，标准方式
// strEncodeWithBase64: 待解码的字节流
func DecodeFromBase64(strEncodeWithBase64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(strEncodeWithBase64)
}

// EncodeToBase64URL 将字节流转换为可视化的base64编码返回，URL兼容方式
// datas: 待编码的字节流
// result: 1 md5值转换为大写；-1 md5值转换为小写；0 md5值不做整理
func EncodeToBase64URL(datas []byte, result int) string {
	base64Result := base64.URLEncoding.EncodeToString(datas)
	if result == 1 {
		return strings.ToUpper(base64Result)
	} else if result == -1 {
		return strings.ToLower(base64Result)
	}
	return base64Result
}

// DecodeFromBase64URL 将base64编码后的字符串解码为原始字节流，URL兼容方式
// strEncodeWithBase64: 待解码的字节流
func DecodeFromBase64URL(strEncodeWithBase64 string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(strEncodeWithBase64)
}
