package encrypt

import (
	"crypto/cipher"
	"fmt"
	chacha "golang.org/x/crypto/chacha20poly1305"
)

const (
	DefaultChachaAad = "pek0rApiratEcospLayGirl1"
	DefaultChachaNonce = "fUbukik0rOneyubiyuBi1919"
)

type Chacha20Poly1305 struct {
	key  []byte 	// 64bit
	Nonce	[]byte  // 24bit
	Aad		[]byte  // additional data
	aead    cipher.AEAD
}

func NewChacha20Poly1305(key ,nonce ,aad string) (*Chacha20Poly1305 ,error) {
	if len(key) != 32 || len(nonce) != 24 {
		return nil ,fmt.Errorf("encrypt.NewChacha20Poly1305 params error(key %s ,nonce %s ,aad %s)" ,key ,nonce ,aad)
	}
	aead ,e := chacha.NewX([]byte(key))
	if e != nil{
		return nil ,e
	}
	return &Chacha20Poly1305{
		key:   []byte(key),
		Nonce: []byte(nonce),
		Aad:   []byte(aad),
		aead:  aead,
	} ,nil
}

func (c *Chacha20Poly1305) Encrypt(plain []byte) ([]byte, error) {
	return c.aead.Seal(nil ,c.Nonce ,plain ,c.Aad) ,nil
}

func (c Chacha20Poly1305) Decrypt(ct []byte) ([]byte ,error) {
	return c.aead.Open(nil ,c.Nonce ,ct ,c.Aad)
}