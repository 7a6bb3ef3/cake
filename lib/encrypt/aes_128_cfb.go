package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/pkg/errors"
	"io"
)

var defaultAes128cfb Encryptor

func SetDefaultAES128CFB(key string) error{
	d ,e := NewAES128CFB(key)
	if e != nil{
		return e
	}
	defaultAes128cfb = d
	return nil
}

type AES128CFB struct {
	key string
	block cipher.Block
}

func (a AES128CFB) Encrypt(in []byte) ([]byte,error) {
	out := make([]byte, aes.BlockSize+len(in))
	iv := out[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(a.block, iv)
	stream.XORKeyStream(out[aes.BlockSize:], in)
	return out ,nil
}

func (a AES128CFB) Decrypt(in []byte) (out []byte, err error) {
	if len(in) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := in[:aes.BlockSize]
	in = in[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(a.block, iv)
	stream.XORKeyStream(in, in)
	return in ,nil
}

func NewAES128CFB(key string) (Encryptor ,error){
	if len(key) != 16 {
		return AES128CFB{} ,errors.Errorf("length of key must be 16")
	}
	block, e := aes.NewCipher([]byte(key))
	if e != nil {
		return AES128CFB{} ,e
	}
	return AES128CFB{
		key: key,
		block: block,
	}, nil
}