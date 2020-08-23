package pool

import (
	"sync"
)

func NewBufPool(bufsize  int) *BufferPool{
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{}{
				return make([]byte ,bufsize)
			},
		},
	}
}

type BufferPool struct {
	pool 	 sync.Pool
}

func (p *BufferPool) Get() []byte{
	return p.pool.Get().([]byte)
}

func (p *BufferPool) Put(buf []byte){
	p.pool.Put(buf)
}