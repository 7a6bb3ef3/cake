package main

import (
	"errors"
	"net"
	"sync"

	"github.com/nynicg/cake/lib/log"

)


func NewTcpConnPool(maxConnnum int) *TcpConnPool{
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

func (p *TcpConnPool)GetLocalTcpConn() (net.Conn ,error){
	select{
	case p.localTks <- struct{}{}:
		log.Debug("get conn -> " ,len(p.localTks))
		conn ,ok := p.localTcpPool.Get().(*net.TCPConn)
		if !ok {
			return p.localTcpPool.New().(net.Conn) ,nil
		}
		return conn ,nil
	default:
		return nil ,errors.New("too many connection ,wait for next available ticket")
	}
}

func (p *TcpConnPool)FreeLocalTcpConn(conn net.Conn) {
	log.Debug("free conn -> " ,len(p.localTks))
	_ = <- p.localTks
	conn.Close()
	p.localTcpPool.Put(conn)
}
