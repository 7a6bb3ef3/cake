package cryptor

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
)

type HMAC struct {
	hash hash.Hash
}

func NewHMAC(key string) *HMAC{
	return &HMAC{hash: hmac.New(sha256.New ,[]byte(key))}
}

func (h *HMAC) SumN(in []byte ,n int) ([]byte ,error){
	out := h.hash.Sum(in)
	if n < 0 || n > len(out) {
		return nil ,fmt.Errorf("out of bounds N %d ,expected [0 ,%d]" ,n ,len(out))
	}
	return out[:n] ,nil
}

func (h *HMAC) Sum(in []byte) []byte{
	return h.hash.Sum(in)
}

func (h *HMAC) SumAhoyHandshake(cmd byte,uid string,length int) ([]byte ,error){
	suminput := append([]byte{cmd ,cmd} ,[]byte(uid)...)
	return h.SumN(suminput ,length)
}
