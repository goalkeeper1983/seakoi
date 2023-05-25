package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"math/big"
	"os"
	"strings"

	"go.uber.org/zap"
)

func GenerateRsaPKCS1KeyToMemory() ([]byte, []byte) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(publicKey)}
	publicKeyByte := pem.EncodeToMemory(block)

	block = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	privateKeyByte := pem.EncodeToMemory(block)

	return publicKeyByte, privateKeyByte
}

func GenerateRsaPKCS8KeyToMemory() ([]byte, []byte) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey

	publicKeyByte, _ := x509.MarshalPKIXPublicKey(publicKey)
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyByte}
	publicKeyByte = pem.EncodeToMemory(block)

	privateKeyByte, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	block = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyByte}
	privateKeyByte = pem.EncodeToMemory(block)

	return publicKeyByte, privateKeyByte
}

func GenerateRsaPKCS1KeyToFile(dir string) {
	for i := 0; i != 100; i++ {
		Log.Info(RunFuncName(), zap.Any("index", i))
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		publicKey := &privateKey.PublicKey

		block := &pem.Block{Type: "PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(publicKey)}
		file, err := os.Create(fmt.Sprintf(dir+"/pkcs1/PKCS1_%03d.pub", i))
		if err != nil {
			Log.Panic(err.Error())
		}
		err = pem.Encode(file, block)
		if err != nil {
			Log.Panic(err.Error())
		}
		file.Close()

		block = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
		file, err = os.Create(fmt.Sprintf(dir+"/pkcs1/PKCS1_%03d.key", i))
		if err != nil {
			Log.Panic(err.Error())
		}
		err = pem.Encode(file, block)
		if err != nil {
			Log.Panic(err.Error())
		}
		file.Close()
	}
}

func GenerateRsaPKCS8KeyToFile(dir string) {
	for i := 0; i != 100; i++ {
		Log.Info(RunFuncName(), zap.Any("index", i))
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		publicKey := &privateKey.PublicKey

		publicKeyByte, _ := x509.MarshalPKIXPublicKey(publicKey)
		block := &pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyByte}
		file, err := os.Create(fmt.Sprintf(dir+"/pkcs8/PKCS8_%03d.pub", i))
		if err != nil {
			Log.Panic(err.Error())
		}
		err = pem.Encode(file, block)
		if err != nil {
			Log.Panic(err.Error())
		}
		file.Close()

		privateKeyByte, _ := x509.MarshalPKCS8PrivateKey(privateKey)
		block = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyByte}
		file, err = os.Create(fmt.Sprintf(dir+"/pkcs8/PKCS8_%03d.key", i))
		if err != nil {
			Log.Panic(err.Error())
		}
		err = pem.Encode(file, block)
		if err != nil {
			Log.Panic(err.Error())
		}
		file.Close()
	}
}

type Rsa struct {
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
}

func (This *Rsa) GetPublicKey() *rsa.PublicKey {
	if This.rsaPublicKey != nil {
		return This.rsaPublicKey
	}
	return nil
}

func (This *Rsa) GetPrivateKey() *rsa.PrivateKey {
	if This.rsaPrivateKey != nil {
		return This.rsaPrivateKey
	}
	return nil
}

// InitRsaPublicKeyPKCS1FromMemory []byte初始化Pkcs1公钥
func (This *Rsa) InitRsaPublicKeyPKCS1FromMemory(pub []byte) {
	var err error
	block, _ := pem.Decode(pub)
	This.rsaPublicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		Log.Panic(err.Error())
	}
}

// InitRsaPrivateKeyPKCS1FromMemory []byte初始化Pkcs1私钥
func (This *Rsa) InitRsaPrivateKeyPKCS1FromMemory(pri []byte) {
	var err error
	block, _ := pem.Decode(pri)
	This.rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		Log.Panic(err.Error())
	}
}

// InitRsaPublicKeyPKCS1 从文件初始化PKCS1公钥
func (This *Rsa) InitRsaPublicKeyPKCS1(pub string) {
	var err error
	buf := ReadFileBytes(pub)
	block, _ := pem.Decode(buf)

	This.rsaPublicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		Log.Panic(err.Error())
	}
}

// InitRsaPrivateKeyPKCS1 从文件初始化PKCS1私钥
func (This *Rsa) InitRsaPrivateKeyPKCS1(pri string) {
	var err error
	buf := ReadFileBytes(pri)
	block, _ := pem.Decode(buf)
	This.rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		Log.Panic(err.Error())
	}
}

// InitRsaPublicKeyPKCS8FromMemory []byte初始化Pkcs8公钥
func (This *Rsa) InitRsaPublicKeyPKCS8FromMemory(pub []byte) {
	block, _ := pem.Decode(pub)
	Any, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		Log.Panic(err.Error())
	}
	This.rsaPublicKey = Any.(*rsa.PublicKey)
}

// InitRsaPrivateKeyPKCS8FromMemory []byte初始化Pkcs8私钥
func (This *Rsa) InitRsaPrivateKeyPKCS8FromMemory(pri []byte) {
	block, _ := pem.Decode(pri)
	Any, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		Log.Panic(err.Error())
	}
	This.rsaPrivateKey = Any.(*rsa.PrivateKey)
}

// InitRsaPublicKeyPKCS8 从文件初始化PKCS8公钥
func (This *Rsa) InitRsaPublicKeyPKCS8(pub string) {
	buf := ReadFileBytes(pub)
	block, _ := pem.Decode(buf)
	Any, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		Log.Panic(err.Error())
	}
	This.rsaPublicKey = Any.(*rsa.PublicKey)
}

// InitRsaPrivateKeyPKCS8 从文件初始化Pkcs8私钥
func (This *Rsa) InitRsaPrivateKeyPKCS8(pri string) {
	buf := ReadFileBytes(pri)
	block, _ := pem.Decode(buf)
	Any, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		Log.Panic(err.Error())
	}
	This.rsaPrivateKey = Any.(*rsa.PrivateKey)
}

func (This *Rsa) InitRsaFromXmlData(publicKeyXmlData, privateKeyXmlData []byte) {
	b64d := func(str string) []byte {
		str = strings.TrimSpace(str)
		decoded := Base64Decode([]byte(str))
		return decoded
	}

	b64bigint := func(str string) *big.Int {
		bint := &big.Int{}
		bint.SetBytes(b64d(str))
		return bint
	}

	xrk := new(struct {
		Modulus  string
		Exponent string
		P        string
		Q        string
		DP       string
		DQ       string
		InverseQ string
		D        string
	})

	if privateKeyXmlData != nil {
		err := xml.Unmarshal(privateKeyXmlData, &xrk)
		if err != nil {
			Log.Panic(err.Error())
		}
		This.rsaPrivateKey = &rsa.PrivateKey{
			PublicKey: rsa.PublicKey{
				N: b64bigint(xrk.Modulus),
				E: int(b64bigint(xrk.Exponent).Int64()),
			},
			D:      b64bigint(xrk.D),
			Primes: []*big.Int{b64bigint(xrk.P), b64bigint(xrk.Q)},
			Precomputed: rsa.PrecomputedValues{
				Dp:        b64bigint(xrk.DP),
				Dq:        b64bigint(xrk.DQ),
				Qinv:      b64bigint(xrk.InverseQ),
				CRTValues: ([]rsa.CRTValue)(nil),
			},
		}
		This.rsaPublicKey = &This.rsaPrivateKey.PublicKey
		return
	}

	if publicKeyXmlData != nil {
		err := xml.Unmarshal(publicKeyXmlData, &xrk)
		if err != nil {
			Log.Panic(err.Error())
		}
		This.rsaPublicKey = &rsa.PublicKey{
			N: b64bigint(xrk.Modulus),
			E: int(b64bigint(xrk.Exponent).Int64()),
		}
	}
}

// RsaPublicKeyEncrypt 公钥加密
func (This *Rsa) RsaPublicKeyEncrypt(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, This.rsaPublicKey, data)
}

// RsaPrivateKeyDecrypt 私钥解密
func (This *Rsa) RsaPrivateKeyDecrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, This.rsaPrivateKey, data)
}

// RsaPublicKeyEncryptByBase64 公钥加密后转base64
func (This *Rsa) RsaPublicKeyEncryptByBase64(data []byte) ([]byte, error) {
	data, err := rsa.EncryptPKCS1v15(rand.Reader, This.rsaPublicKey, data)
	if err != nil {
		return nil, err
	}
	data = Base64Encode(data)
	return data, nil
}

// RsaPrivateKeyDecryptByBase64 私钥解密base64
func (This *Rsa) RsaPrivateKeyDecryptByBase64(base64Data []byte) ([]byte, error) {
	data := Base64Decode(base64Data)
	return rsa.DecryptPKCS1v15(rand.Reader, This.rsaPrivateKey, data)
}
