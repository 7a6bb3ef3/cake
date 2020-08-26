package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type AES128GCM struct {
	nonce	[]byte
	block	cipher.Block
	aead	cipher.AEAD
}

// non-hex key and nonce
func NewAES128GCM(key string) (Cryptor ,error){
	nonce := sha256N(key ,12)
	if len(key) != 32 || len(nonce) != 12 {
		return nil ,errors.New("unexpected length of key or nonce")
	}
	block ,e := aes.NewCipher([]byte(key))
	if e != nil{
		return nil ,e
	}
	ae ,e := cipher.NewGCM(block)
	if e != nil{
		return nil ,e
	}
	re := &AES128GCM{
		nonce: nonce,
		block: block,
		aead:  ae,
	}
	return re ,nil
}

func (a AES128GCM) Encrypt(in []byte) (out []byte, err error) {
	return a.aead.Seal(nil ,a.nonce ,in ,nil) ,nil
}

func (a AES128GCM) Decrypt(in []byte) (out []byte, err error) {
	return a.aead.Open(nil ,a.nonce ,in ,nil)
}

