package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
)

type AES128CBC struct {
	key		[]byte
	block 	cipher.Block
	iv      []byte
}

func NewAES128CBC(key ,iv string) (StreamEncryptor ,error){
	if len(key) != 16 || len(iv) != 16 {
		return AES128CBC{} ,errors.New("the length of key and iv must be 16")
	}
	ae := AES128CBC{
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

func (a AES128CBC) EncryptStream(out io.Writer ,in io.Reader) error {
	buf := bytes.NewBuffer([]byte{})
	if _, e := io.Copy(buf ,in);e != nil{
		return e
	}
	encdata ,e := a.Encrypt(buf.Bytes())
	if e != nil{
		return e
	}
	_ ,e = out.Write(encdata)
	return e
}

func (a AES128CBC) DecryptStream(out io.Writer ,in io.Reader) error {
	buf := bytes.NewBuffer([]byte{})
	if _, e := io.Copy(buf ,in);e != nil{
		return e
	}
	plainData ,e := a.Decrypt(buf.Bytes())
	if e != nil{
		return e
	}
	_ ,e = out.Write(plainData)
	return e
}

func (a AES128CBC) Encrypt(in []byte) (out []byte, err error) {

	blockSize := a.block.BlockSize()
	in = a.pkcs7Padding(in, blockSize)

	blockMode := cipher.NewCBCEncrypter(a.block, a.iv)
	crypted := make([]byte, len(in))
	blockMode.CryptBlocks(crypted, in)

	return crypted, nil
}

func (a AES128CBC) Decrypt(in []byte) (out []byte, err error) {
	if len(in) == 0 {
		return nil ,errors.New("empty input data to be decrypted")
	}
	blockMode := cipher.NewCBCDecrypter(a.block, a.iv)
	origData := make([]byte, len(in))

	blockMode.CryptBlocks(origData, in)
	return a.pkcs7UnPadding(origData)
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
	if i < 0 || i > len(origData) {
		return nil ,errors.New("aes decrypt failed. slice bounds out of range")
	}
	return origData[:(length - unpadding)] ,nil
}