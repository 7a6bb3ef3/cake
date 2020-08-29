package main

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"
)

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

func TestUIDManager(t *testing.T) {
	uid1 := randStr()[:32]
	uid2 := randStr()[:32]

	go func() {
		tk := time.NewTicker(time.Second * 8)
		for range tk.C {
			DefUidManager().refreshCache()
		}
	}()

	go func() {
		for {
			DefUidManager().VerifyHMAC([]byte("error hmac"))
			time.Sleep(time.Millisecond * 30)
		}
	}()

	DefUidManager().RegisterUid(uid1, UIDInfo{})
	DefUidManager().RegisterUid(uid2, UIDInfo{})
	time.Sleep(time.Second * 30)
}

// 101163 ns/op	   64854 B/op	    1042 allocs/op
func BenchmarkUIDManager(b *testing.B) {
	DefUidManager().RegisterUid(randStr()[:32], UIDInfo{})
	DefUidManager().RegisterUid(randStr()[:32], UIDInfo{})
	DefUidManager().RegisterUid(randStr()[:32], UIDInfo{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		DefUidManager().refreshCache()
	}
}
