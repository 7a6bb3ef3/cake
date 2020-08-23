package pool

import (
	"net"
	"sync"

	"github.com/nynicg/cake/lib/log"
)


func NewTcpConnPool(maxConnnum int) *TcpConnPool {
	p := &TcpConnPool{}
	p.localTks = make(chan struct{} ,maxConnnum)
	p.localTcpPool = sync.Pool{
		New: func() interface{}{
			return &net.TCPConn{}
		},
	}
	return p
}

type TcpConnPool struct {
	localTks chan struct{}
	localTcpPool sync.Pool
}

func (p *TcpConnPool)GetLocalTcpConn() net.Conn{
	select{
	case p.localTks <- struct{}{}:
		log.Debug("get conn -> " ,len(p.localTks))
		conn ,ok := p.localTcpPool.Get().(*net.TCPConn)
		if !ok {
			return p.localTcpPool.New().(net.Conn)
		}
		return conn
	}
}

func (p *TcpConnPool)FreeLocalTcpConn(conn net.Conn) {
	_ = <- p.localTks
	conn.Close()
	p.localTcpPool.Put(conn)
	log.Debug("free conn -> " ,len(p.localTks))
}
