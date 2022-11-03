package util

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io/ioutil"
)

/*
===  数据的加密解密
*/

const (
	PKCS5_PADDING = "PKCS5"
	PKCS7_PADDING = "PKCS7"
	ZEROS_PADDING = "ZEROS"
)

//AesCBCEncrypt 加密
func AesCBCEncrypt(src, key []byte, paddingMode string) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := cipherBlock.BlockSize()

	src = padding(src, paddingMode, blockSize)
	//加密
	encryptData := make([]byte, len(src))

	mode := cipher.NewCBCEncrypter(cipherBlock, key[:blockSize])
	mode.CryptBlocks(encryptData, src)

	return encryptData, nil
}

//AesCBCDecrypt 解密
func AesCBCDecrypt(src, key []byte, paddingMode string) ([]byte, error) {
	//未知  后面等测试
	//decodeString, err := hex.DecodeString(string(src))
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, err
	//}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := cipherBlock.BlockSize()

	dst := make([]byte, len(src))

	mode := cipher.NewCBCDecrypter(cipherBlock, key[:blockSize])
	mode.CryptBlocks(dst, src)

	return unPadding(dst, paddingMode), nil
}

func padding(src []byte, paddingMode string, blockSize int) []byte {
	switch paddingMode {
	case PKCS5_PADDING:

	case PKCS7_PADDING:
		src = pkcs7Padding(src, blockSize)
	case ZEROS_PADDING:

	}
	return src
}

func unPadding(src []byte, paddingMode string) []byte {
	var res []byte
	switch paddingMode {
	case PKCS5_PADDING:

	case PKCS7_PADDING:
		res = pkcs7UnPadding(src)
	case ZEROS_PADDING:

	}
	return res
}

//pkcs7Padding  补码
func pkcs7Padding(src []byte, blockSize int) []byte {
	paddingLen := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)
	return append(src, padtext...)
}

//pkcs7UnPadding 去码
func pkcs7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

/*
==============数据的解析处理
*/

func Zip(data []byte) ([]byte, error) {
	var (
		err  error
		gzwr *gzip.Writer
	)
	//底层 io.writer 接口
	var b bytes.Buffer
	if gzwr, err = gzip.NewWriterLevel(&b, gzip.BestCompression); err != nil {
		return nil, err
	}
	//将 data 数据压缩后 写入下层的 io.writer 接口 也就是b
	if _, err = gzwr.Write(data); err != nil {
		return nil, err
	}
	//将缓冲中的压缩数据全部刷新到底层 的io.writer 接口中
	if err = gzwr.Flush(); err != nil {
		return nil, err
	}
	if err = gzwr.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func UnZip(data []byte) ([]byte, error) {
	var (
		err       error
		reader    *gzip.Reader
		unZipData []byte
	)
	b := new(bytes.Buffer)
	if err = binary.Write(b, binary.LittleEndian, data); err != nil {
		return nil, err
	}
	if reader, err = gzip.NewReader(b); err != nil {
		return nil, err
	}
	defer reader.Close()

	if unZipData, err = ioutil.ReadAll(reader); err != nil {
		return nil, err
	}
	return unZipData, nil
}
