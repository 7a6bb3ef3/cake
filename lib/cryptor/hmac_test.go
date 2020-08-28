package cryptor

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestHMAC(t *testing.T) {
	in := buildInput()
	for _ ,v := range in{
		b := HMAC(v)
		t.Log(b)
	}
	t.Log(strconv.Itoa(int(time.Now().Unix())))
}

func TestVerifyHMAC(t *testing.T) {
	uid := "12323121232312123231212323121235"
	hs := HMAC(uid)
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	t.Log(VerifyHMAC(uid ,hs))

	t.Log(HMACAllTime(uid))
	t.Log(hs)
}

// 1048 ns/op	     640 B/op	      11 allocs/op
func BenchmarkHMAC(b *testing.B) {
	uid := "12323121232312123231212323121235"
	b.ReportAllocs()
	for i:=0;i<b.N;i++{
		HMAC(uid)
	}
}
