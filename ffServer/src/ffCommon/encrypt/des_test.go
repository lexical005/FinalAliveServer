package encrypt

import (
	"ffCommon/log/log"

	"encoding/base64"
	"testing"
)

func testDes() {
	key := []byte("sfe023f_")
	result, err := DesEncryptPaddingPKCS5([]byte("polaris@studygolang"), key)
	if err != nil {
		panic(err)
	}
	log.RunLogger.Println(base64.StdEncoding.EncodeToString(result))
	origData, err := DesDecryptPaddingPKCS5(result, key)
	if err != nil {
		panic(err)
	}
	log.RunLogger.Println(string(origData))

	result, err = DesEncryptPaddingZero([]byte("polaris@studygolang"), key)
	if err != nil {
		panic(err)
	}
	log.RunLogger.Println(base64.StdEncoding.EncodeToString(result))
	origData, err = DesEncryptPaddingZero(result, key)
	if err != nil {
		panic(err)
	}
	log.RunLogger.Println(string(origData))
}

func test3Des() {
	key := []byte("sfe023f_sefiel#fi32lf3e!")
	result, err := TripleDesEncryptPaddingPKCS5([]byte("polaris@studygol"), key)
	if err != nil {
		panic(err)
	}
	log.RunLogger.Println(base64.StdEncoding.EncodeToString(result))
	origData, err := TripleDesDecryptPaddingPKCS5(result, key)
	if err != nil {
		panic(err)
	}
	log.RunLogger.Println(string(origData))

	key = []byte("sfe023f_sefiel#fi32lf3e!")
	result, err = TripleDesEncryptPaddingZero([]byte("polaris@studygol"), key)
	if err != nil {
		panic(err)
	}
	log.RunLogger.Println(base64.StdEncoding.EncodeToString(result))
	origData, err = TripleDesEncryptPaddingZero(result, key)
	if err != nil {
		panic(err)
	}
	log.RunLogger.Println(string(origData))
}

func Test_Des(t *testing.T) {
	testDes()
	test3Des()
}
