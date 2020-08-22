package main

import (
	"fmt"
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/socks5"
	"io"
	"net"
	"time"
)

func startLocalSocksProxy() net.Conn{
	log.Info("Socks5 listen on " ,config.LocalSocksAddr)
	ls ,e := net.Listen("tcp" ,config.LocalSocksAddr)
	if e != nil{
		log.Panic(e)
	}
	pool := NewTcpConnPool(config.MaxLocalConnNum)
	for {
		cliconn ,e := pool.GetLocalTcpConn()
		if e != nil{
			log.Errorx("try to get local tcp conn from pool " ,e)
			time.Sleep(time.Millisecond * 64)
			continue
		}
		cliconn ,e = ls.Accept()
		if e != nil{
			cliconn.Close()
			log.Error(e)
			continue
		}
		go handleCliConn(cliconn ,pool)
	}
}

func handleCliConn(cliconn net.Conn ,pool *TcpConnPool){
	defer func() {
		cliconn.Close()
		pool.FreeLocalTcpConn(cliconn)
	}()
	cliconn.(*net.TCPConn).SetKeepAlive(false)
	if e := socks5.Handshake(cliconn);e != nil{
		log.Errorx("handshake with "+cliconn.RemoteAddr().String() ,e)
		return
	}
	addr, e := socks5.ParseCMD(cliconn)
	if e != nil{
		log.Errorx("parse client cmd and addr" ,e)
		return
	}
	log.Debug("parse remote host -> " ,addr.Address())
	var remote net.Conn
	if Bypass(addr.Host()) {
		remote ,e = net.Dial("tcp" ,addr.Address())
		if e != nil{
			log.Errorx("dail bypassed remote addr " ,e)
			socks5.ProxyFailed(socks5.SocksRespHostUnreachable ,cliconn)
			return
		}
		remote.(*net.TCPConn).SetKeepAlive(false)
	}else{
		// connects to remote socks5 proxy
	}
	defer remote.Close()
	if e := socks5.ProxyOK(cliconn);e != nil{
		log.Errorx("local socks5 sent OK resp " ,e)
		return
	}

	log.Debug("finish handshake ,ready to transport")
	go func() {
		upN ,e := io.Copy(remote ,cliconn)
		if e != nil{
			log.Errorx("copy client request to remote server error " ,e)
			return
		}
		log.Debug(fmt.Sprintf("%s %d bit ↑" ,remote.RemoteAddr() ,upN))
	}()
	downN ,e := io.Copy(cliconn ,remote)
	if e != nil{
		log.Errorx("copy server response to client error " ,e)
		return
	}
	log.Debug(fmt.Sprintf("%s %d bit ↓" ,remote.RemoteAddr() ,downN))
}
