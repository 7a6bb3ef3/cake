package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var defaultAES128CBC Encryptor

type AES128CBC struct {
	key		[]byte
	block 	cipher.Block
	iv      []byte
}

func SetDefaultAES128CBC(key ,iv string) error{
	en ,e :=  NewAES128CBC(key ,iv)
	defaultAES128CBC = en
	return e
}

func NewAES128CBC(key ,iv string) (Encryptor ,error){
	if len(key) != 16 || len(iv) != 16 {
		return &AES128CBC{} ,errors.New("the length of key and iv must be 16")
	}
	ae := &AES128CBC{
		key: []byte(key),
		iv:  []byte(iv),
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ae, err
	}
	ae.block = block
	return ae ,nil
}

func (a *AES128CBC) Encrypt(in []byte) (out []byte, err error) {

	return _testReverse(in) ,nil
	//blockSize := a.block.BlockSize()
	//in = a.pkcs7Padding(in, blockSize)
	//
	//blockMode := cipher.NewCBCEncrypter(a.block, a.iv)
	//crypted := make([]byte, len(in))
	//blockMode.CryptBlocks(crypted, in)
	//return crypted, nil
}

func (a *AES128CBC) Decrypt(in []byte) (out []byte, err error) {
	return _testReverse(in) ,nil
	//if len(in) == 0 {
	//	return in ,nil
	//}
	//if len(in) % a.block.BlockSize() != 0 {
	//	return in ,errors.New("input data seems not like an encrypted byte slice")
	//}
	//blockMode := cipher.NewCBCDecrypter(a.block, a.iv)
	//origData := make([]byte, len(in))
	//blockMode.CryptBlocks(origData, in)
	//return a.pkcs7UnPadding(origData)
}

func (a AES128CBC) pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (a AES128CBC)pkcs7UnPadding(origData []byte) ([]byte ,error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	i := length - unpadding
	if i < 0 {
		return nil ,errors.New("aes decrypt failed. slice bounds out of range")
	}
	return origData[:(length - unpadding)] ,nil
}


func _testReverse(in []byte) []byte{
	length := len(in)
	out := make([]byte ,length)
	for i:=0;i<length;i++{
		out[length - 1 - i] = in[i]
	}
	return out
}