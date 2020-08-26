package cryptor

import (
	"encoding/hex"
	"testing"
)

func TestNewAES128GCM(t *testing.T) {
	key := "808182838485868788898a8b8c8d1233"
	gcm ,e := NewAES128GCM(key)
	if e != nil{
		t.Fatal(e)
	}

	gcmde ,e := NewAES128GCM(key)
	if e != nil{
		t.Fatal(e)
	}
	cases := buildInput()
	for _ ,v := range cases{
		inb ,e := hex.DecodeString(v)
		if e != nil{
			t.Errorf("decode %s %s" ,v ,e)
			continue
		}
		out ,e := gcm.Encrypt(inb)
		t.Logf("len in %d ,len out %d" ,len(inb) ,len(out))
		if e != nil{
			t.Errorf("cryptor %s %s" ,v ,e)
			continue
		}
		plain ,e := gcmde.Decrypt(out)
		if e != nil{
			t.Errorf("decrypt %s %s" ,v ,e)
			continue
		}
		match := v == hex.EncodeToString(plain)
		//t.Logf("%t %s -> %v -> %s" ,match ,v ,out ,plain)
		t.Log(match)
	}
}
