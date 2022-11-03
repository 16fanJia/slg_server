package util

import (
	"fmt"
	"testing"
)

func TestAesCBCEncrypt(t *testing.T) {
	src := "加密功能测试"
	key := "123456781234567812345678"
	fmt.Println("原文：", src)

	encryptCode, _ := AesCBCEncrypt([]byte(src), []byte(key), PKCS7_PADDING)
	fmt.Println("密文：", encryptCode)

	decryptCode, _ := AesCBCDecrypt(encryptCode, []byte(key), PKCS7_PADDING)
	fmt.Println("解密结果：", string(decryptCode))
}
