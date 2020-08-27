package cryptor

import (
	"encoding/hex"
	"math/rand"
	"testing"
)

func TestChacha20Poly1305_Encrypt(t *testing.T) {
	key := "808182838485868788898a8b8c8d1233"
	cha, e := NewChacha20Poly1305(key)
	if e != nil {
		t.Fatal(e)
	}
	cases := buildInput()
	for _, v := range cases {
		inb, e := hex.DecodeString(v)
		if e != nil {
			t.Errorf("decode %s %s", v, e)
			continue
		}
		out, e := cha.Encrypt(inb)
		t.Logf("len in %d ,len out %d", len(inb), len(out))
		if e != nil {
			t.Errorf("cryptor %s %s", v, e)
			continue
		}
		plain, e := cha.Decrypt(out)
		if e != nil {
			t.Errorf("decrypt %s %s", v, e)
			continue
		}
		match := v == hex.EncodeToString(plain)
		//t.Logf("%t %s -> %v -> %s" ,match ,v ,out ,plain)
		t.Log(match)
	}

}

func buildInput() []string {
	l := 20
	var out = make([]string, l)
	for i := 0; i < l; i++ {
		out[i] = randStr()
	}
	return out
}

func randStr() string {
	bits := rand.Intn(1000) + 500
	byts := make([]byte, bits)
	for i := 0; i < bits; i++ {
		byts[i] = byte(rand.Intn(256))
	}
	return hex.EncodeToString(byts)
}
