package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/nynicg/cake/lib/ahoy"
	"github.com/nynicg/cake/lib/cryptor"
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/pool"
	"github.com/nynicg/cake/lib/socks5"
)

var bufpool *pool.BufferPool

func init(){
	bufpool = pool.NewBufPool(32 * 1024)
}

func startLocalSocksProxy(encryptor cryptor.Cryptor){
	ls ,e := net.Listen("tcp" ,config.LocalSocksAddr)
	if e != nil{
		log.Panic(e)
	}
	log.Info("Socks5 listen on " ,config.LocalSocksAddr)
	pl := pool.NewTcpConnPool(config.MaxLocalConnNum)
	for {
		cliconn := pl.GetLocalTcpConn()
		cliconn ,e = ls.Accept()
		if e != nil{
			log.Errorx("accept new client conn " ,e)
			continue
		}
		go handleCliConn(cliconn ,pl ,encryptor)
	}
}

func handleCliConn(cliconn net.Conn ,pl *pool.TcpConnPool,encryptor cryptor.Cryptor){
	defer func() {
		cliconn.Close()
		pl.FreeLocalTcpConn(cliconn)
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

	var remote net.Conn
	var bypass int
	if i ,ok := GetDomainCache(addr.Host());ok{
		bypass = i
	}else{
		bypass = Bypass(addr.Host())
		PutDomainCache(addr.Host() ,bypass)
	}
	var cryptorSelect cryptor.Cryptor
	switch bypass {
	case BypassDiscard:
		socks5.ProxyFailed(socks5.SocksRespHostUnreachable ,cliconn)
		return
	case BypassTrue:
		cryptorSelect = cryptor.GetTypePlain()
		remote ,e = net.Dial("tcp" ,addr.Address())
		if e != nil{
			log.Errorx("dail bypassed remote addr " ,e)
			socks5.ProxyFailed(socks5.SocksRespHostUnreachable ,cliconn)
			return
		}
	case BypassProxy:
		cryptorSelect = encryptor
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
	default:
		socks5.ProxyFailed(socks5.SocksRespServErr ,cliconn)
		return
	}
	defer remote.Close()
	remote.(*net.TCPConn).SetKeepAlive(false)
	if e := socks5.ProxyOK(cliconn);e != nil{
		log.Errorx("local socks5 sent OK resp " ,e)
		return
	}

	outboundEnv := &ahoy.CopyEnv{
		ReaderWithLength: false,
		WriterNeedLength: true,
		CryptFunc:        cryptorSelect.Encrypt,
		BufPool:          bufpool,
		Bypass:			  bypass == BypassTrue,
	}
	inboundEnv := &ahoy.CopyEnv{
		ReaderWithLength: true,
		WriterNeedLength: false,
		CryptFunc:        cryptorSelect.Decrypt,
		BufPool:          bufpool,
		Bypass:			  bypass == BypassTrue,
	}

	var(
		up int
		down int
	)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		up ,e = ahoy.CopyConn(remote ,cliconn ,outboundEnv)
		if e != nil{
			log.Warn("proxy request." ,e)
			return
		}
	}()
	go func() {
		defer wg.Done()
		down ,e = ahoy.CopyConn(cliconn ,remote ,inboundEnv)
		if e != nil{
			log.Warn("server resp." ,e)
			return
		}
	}()
	wg.Wait()
	log.Info(fmt.Sprintf("%s ,%d ↑ ,%d ↓ bytes" ,addr.Host() ,up ,down))
}


func handshakeRemote(remote net.Conn ,proxyhost string) error{
	if len(proxyhost) > 255 {
		return errors.New("host addr is too long(>255)")
	}
	index ,e := cryptor.GetStreamEncryptorIndexByName(config.EncryptType)
	if e != nil{
		return e
	}
	req := ahoy.RemoteConnectRequest{
		Encryption: 	byte(index),
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
	return nil
}
