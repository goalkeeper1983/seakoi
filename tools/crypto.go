package tools

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"io"
	"log"
)

// MD5 加密
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ZlibUnzip(src []byte) []byte {
	b := bytes.NewReader(src)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	if _, err := io.Copy(&out, r); err != nil {
		log.Panicln(err.Error())
	}
	return out.Bytes()
}

func ZlibZip(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	if _, err := w.Write(src); err != nil {
		log.Panicln(err.Error())
	}
	err := w.Close()
	if err != nil {
		return nil
	}
	return in.Bytes()
}

func Base64Encode(data []byte) []byte {
	s := base64.StdEncoding.EncodeToString(data)
	return []byte(s)
}

func Base64Decode(data []byte) []byte {
	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		Log.Panic(RunFuncName(), zap.Any("err", err.Error()))
	}
	return data
}

///////aes/cbc/pkcs7
//加密过程：
//  1、处理数据，对数据进行填充，采用PKCS7（当密钥长度不够时，缺几位补几个几）的方式。
//  2、对数据进行加密，采用AES加密方法中CBC加密模式
//  3、对得到的加密数据，进行base64加密，得到字符串
// 解密过程相反

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

// AesEncrypt 加密
func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	//创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil

}

// AesDecrypt 解密
func AesDecrypt(data []byte, key []byte) ([]byte, error) {
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

// AesEncryptByBase64 Aes加密后进行base64编码
// key长度为16 对应 AES-128
// key长度为24 对应 AES-192
// key长度为32 对应 AES-256
func AesEncryptByBase64(data, key string) string {
	res, err := AesEncrypt([]byte(data), []byte(key))
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(res)
}

// AesDecryptByBase64 Aes加密后进行base64解码
func AesDecryptByBase64(data, key string) string {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}
	decryptData, err := AesDecrypt(dataByte, []byte(key))
	if err != nil {
		panic(err)
	}
	return string(decryptData)
}
