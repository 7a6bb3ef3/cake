package main

import (
	"errors"
	"fmt"
	"github.com/nynicg/cake/lib/ahoy"
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/socks5"
	"io"
	"net"
	"time"
)

const maxRetryTimes = 10

func startLocalSocksProxy(){
	log.Info("Socks5 listen on " ,config.LocalSocksAddr)
	ls ,e := net.Listen("tcp" ,config.LocalSocksAddr)
	if e != nil{
		log.Panic(e)
	}
	pool := NewTcpConnPool(config.MaxLocalConnNum)
	var tried = 0
	for {
		cliconn ,e := pool.GetLocalTcpConn()
		if e != nil{
			tried++
			if tried > maxRetryTimes{
				log.Error("after attempts, there is still no conn available ,system exit")
				return
			}
			log.Errorx("try to get local tcp conn from pool " ,e)
			time.Sleep(time.Second)
			continue
		}
		tried = 0
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
		remote ,e = net.Dial("tcp" ,config.RemoteExitAddr)
		if e != nil{
			log.Errorx("dail remote exit " ,e)
			socks5.ProxyFailed(socks5.SocksRespServErr ,cliconn)
			return
		}
		if e := handshakeRemote(remote ,addr.Address());e != nil{
			log.Errorx("handshake with remote failed " ,e)
			socks5.ProxyFailed(socks5.SocksRespServErr ,cliconn)
			return
		}
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
			log.Warn("copy client request to remote server error " ,e)
			return
		}
		log.Debug(fmt.Sprintf("%s %d bit ↑" ,remote.RemoteAddr() ,upN))
	}()
	downN ,e := io.Copy(cliconn ,remote)
	if e != nil{
		log.Warn("copy server response to client error " ,e)
		return
	}
	log.Debug(fmt.Sprintf("%s %d bit ↓" ,remote.RemoteAddr() ,downN))
}

func handshakeRemote(remote net.Conn ,proxyhost string) error{
	if len(proxyhost) > 255 {
		return errors.New("host addr is too long(>255)")
	}
	req := ahoy.RemoteConnectRequest{
		Version: 		1,
		Command: 		ahoy.CommandConnect,
		AccessKey:      []byte(config.RemoteAccessKey),
		AddrLength:     byte(len(proxyhost)),
		Addr: 			[]byte(proxyhost),
	}
	bts ,e := req.Bytes()
	if e != nil{
		return e
	}
	if _, e := remote.Write(bts);e != nil{
		return e
	}
	// read ok resp {1,1,4,5,1,4}
	buf := make([]byte ,6)
	if _ ,e := io.ReadFull(remote ,buf);e != nil{
		return e
	}
	if buf[0] == 1 && buf[1] == 1 &&buf[2] == 4 &&buf[3] == 5 &&buf[4] == 1 && buf[5] == 4 {
		return nil
	}
	return errors.New("unknown server resp")
}
