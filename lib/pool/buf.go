package pool

import (
	"sync"
)


func NewBufPool(bufsize int) *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, bufsize)
			},
		},
		size: bufsize,
	}
}

type BufferPool struct {
	pool 	sync.Pool
	size	int
}

func (p *BufferPool) resetByte(b []byte) []byte{
	if p.size != len(b) {
		b = b[:p.size]
	}
	for k := range b{
		b[k] = 0
	}
	return b
}

func (p *BufferPool) Get() []byte {
	b := p.pool.Get().([]byte)
	return p.resetByte(b)
}

func (p *BufferPool) Put(buf []byte) {
	p.pool.Put(buf)
}
