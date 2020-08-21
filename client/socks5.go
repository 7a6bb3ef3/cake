package main

import (
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/socks5"
	"net"
)

func startLocalSocksProxy() net.Conn{
	log.Info("Socks5 listen on " ,config.LocalSocksAddr)
	ls ,e := net.Listen("tcp" ,config.LocalSocksAddr)
	if e != nil{
		log.Panic(e)
	}
	for {
		cliconn ,e := ls.Accept()
		log.Debug("get conn from client " ,cliconn.LocalAddr())
		if e != nil{
			cliconn.Close()
			log.Error(e)
			continue
		}
		go handleCliConn(cliconn)
	}
}

func handleCliConn(cliconn net.Conn){
	defer cliconn.Close()
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
	log.Info("parse remote host -> " ,addr.String())
}
