package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

//AES/ECB/PKCS5

func AesEcbPkcs5Decrypt(crypto, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		Log.Panic(err.Error())
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypto))
	blockMode.CryptBlocks(origData, crypto)
	origData = pkcs5UnPadding(origData)
	return origData
}

func AesEcbPkcs5Encrypt(source, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		Log.Panic(err.Error())
	}
	encrypter := NewECBEncrypter(block)
	content := pkcs5Padding(source, block.BlockSize())
	dst := make([]byte, len(content))
	encrypter.CryptBlocks(dst, content)
	return dst
}

func AesEcbPkcs5Base64Decrypt(crypto, key []byte) []byte {
	crypto = Base64Decode(crypto)
	return AesEcbPkcs5Decrypt(crypto, key)
}

func AesEcbPkcs5Base64Encrypt(crypto, key []byte) []byte {
	crypto = AesEcbPkcs5Encrypt(crypto, key)
	return Base64Encode(crypto)
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, paddingText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
