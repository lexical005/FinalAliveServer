package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// MD5 计算字节流的 md5
// datas: 待计算的字节流
// result: 1 md5值转换为大写；-1 md5值转换为小写；0 md5值不做整理
func MD5(datas []byte, result int) string {
	md5Result := md5.Sum(datas)
	if result == 1 {
		return strings.ToUpper(hex.EncodeToString(md5Result[:]))
	} else if result == -1 {
		return strings.ToLower(hex.EncodeToString(md5Result[:]))
	}
	return hex.EncodeToString(md5Result[:])
}
