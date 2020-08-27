package cryptor

import (
	"bytes"
	"testing"
)

func TestXorStream(t *testing.T) {
	cas := buildInput()
	key := "BAby10nStAGec0atBAby10nStAGec0at"
	for _, v := range cas {
		src := []byte(v)
		t.Log(src)
		dst := make([]byte, len(src))
		XorStream(dst, src, key)
		t.Log(dst)
		XorStream(dst[:10], dst[:10], key)
		XorStream(dst[10:], dst[10:], key)
		t.Log(dst)
		// expect dst[:10] == src[:10]  dst[10:] != src[10:]
		t.Log(bytes.Equal(dst, src))
	}
}

func BenchmarkXorStream(b *testing.B) {
	b.Log(fixedNonce)
	msg := []byte("BAby10nStAGec0atBAby10nStAGec0at") // 36
	key := "BAby10nStAGec0atBAby10nStAGec0at"
	//dst := make([]byte ,len(msg))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst := make([]byte, len(msg))
		XorStream(dst, msg, key)
	}
}
