package cryptor

import (
	"crypto/cipher"
	"fmt"

	chacha "golang.org/x/crypto/chacha20poly1305"
)


type Chacha20Poly1305 struct {
	nonce	[]byte  // 24bit
	aead    cipher.AEAD
}

func NewChacha20Poly1305(key string) (Cryptor ,error) {
	nonce := sha256N(key ,24)
	if len(key) != 32  {
		return nil ,fmt.Errorf("cryptor.NewChacha20Poly1305 params error(key %s ,nonce %s)" ,key ,nonce)
	}
	aead ,e := chacha.NewX([]byte(key))
	if e != nil{
		return nil ,e
	}
	return &Chacha20Poly1305{
		nonce: nonce,
		aead:  aead,
	} ,nil
}

func (c *Chacha20Poly1305) Encrypt(plain []byte) ([]byte, error) {
	return c.aead.Seal(nil ,c.nonce ,plain ,nil) ,nil
}

func (c Chacha20Poly1305) Decrypt(ct []byte) ([]byte ,error) {
	return c.aead.Open(nil ,c.nonce ,ct ,nil)
}