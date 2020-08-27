package cryptor

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

// Deprecated
func SumAhoyHandshake(cmd byte, uid string, length int) ([]byte, error) {
	uidbyte := []byte(uid)
	hs := hmac.New(sha256.New, uidbyte)
	suminput := append([]byte{cmd, cmd}, []byte(uid)...)
	out := hs.Sum(suminput)
	if length < 0 || length > len(out) {
		return nil, fmt.Errorf("out of bounds N %d ,expected [0 ,%d]", length, len(out))
	}
	return out[:length], nil
}
