package pool

import (
	"testing"
)

func TestNewBufPool(t *testing.T) {
	pl := NewBufPool(16)
	for i:=0;i<1000;i++{
		b := pl.Get()
		if len(b) != pl.size{
			t.Fail()
			return
		}
		b[0] = 123
		b = append(b, 1)
		pl.Put(b)
	}
}

func BenchmarkBufferPool_Get(b *testing.B) {
	b.ReportAllocs()
	pl := NewBufPool(1024 * 16)
	for i:=0;i<b.N;i++{
		b := pl.Get()
		pl.Put(b)
	}
}
