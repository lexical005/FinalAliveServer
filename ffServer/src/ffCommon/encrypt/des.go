package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

// DesEncryptPaddingPKCS5 des enctypt with CBC mode and PKCS5 padding
func DesEncryptPaddingPKCS5(dataOriginal, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	dataOriginal = paddingPKCS5(dataOriginal, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	dataEncrypted := dataOriginal
	blockMode.CryptBlocks(dataEncrypted, dataOriginal)
	return dataEncrypted, nil
}

// DesDecryptPaddingPKCS5 des dectypt with CBC mode and PKCS5 padding
func DesDecryptPaddingPKCS5(dataEncrypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, key)
	dataOriginal := dataEncrypted
	blockMode.CryptBlocks(dataOriginal, dataEncrypted)
	dataOriginal = unPaddingPKCS5(dataOriginal)
	return dataOriginal, nil
}

// DesEncryptPaddingZero des enctypt with CBC mode and Zero padding
func DesEncryptPaddingZero(dataOriginal, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	dataOriginal = paddingZero(dataOriginal, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	dataEncrypted := dataOriginal
	blockMode.CryptBlocks(dataEncrypted, dataOriginal)
	return dataEncrypted, nil
}

// DesDecryptPaddingZero des dectypt with CBC mode and Zero padding
func DesDecryptPaddingZero(dataEncrypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	dataOriginal := dataEncrypted
	blockMode.CryptBlocks(dataOriginal, dataEncrypted)
	dataOriginal = unPaddingZero(dataOriginal)
	return dataOriginal, nil
}

// TripleDesEncryptPaddingPKCS5 3des enctypt with CBC mode and PKCS5 padding
func TripleDesEncryptPaddingPKCS5(dataOriginal, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	dataOriginal = paddingPKCS5(dataOriginal, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	dataEncrypted := dataOriginal
	blockMode.CryptBlocks(dataEncrypted, dataOriginal)
	return dataEncrypted, nil
}

// TripleDesDecryptPaddingPKCS5 3des dectypt with CBC mode and PKCS5 padding
func TripleDesDecryptPaddingPKCS5(dataEncrypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	dataOriginal := dataEncrypted
	blockMode.CryptBlocks(dataOriginal, dataEncrypted)
	dataOriginal = unPaddingPKCS5(dataOriginal)
	return dataOriginal, nil
}

// TripleDesEncryptPaddingZero 3des enctypt with CBC mode and Zero padding
func TripleDesEncryptPaddingZero(dataOriginal, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	dataOriginal = paddingZero(dataOriginal, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	dataEncrypted := dataOriginal
	blockMode.CryptBlocks(dataEncrypted, dataOriginal)
	return dataEncrypted, nil
}

// TripleDesDecryptPaddingZero 3des dectypt with CBC mode and Zero padding
func TripleDesDecryptPaddingZero(dataEncrypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	dataOriginal := dataEncrypted
	blockMode.CryptBlocks(dataOriginal, dataEncrypted)
	dataOriginal = unPaddingZero(dataOriginal)
	return dataOriginal, nil
}

func paddingZero(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func unPaddingZero(dataOriginal []byte) []byte {
	return bytes.TrimRightFunc(dataOriginal, func(r rune) bool {
		return r == rune(0)
	})
}

func paddingPKCS5(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func unPaddingPKCS5(dataOriginal []byte) []byte {
	length := len(dataOriginal)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(dataOriginal[length-1])
	return dataOriginal[:(length - unpadding)]
}
