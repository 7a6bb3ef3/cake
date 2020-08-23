package main

import (
	"errors"
	"fmt"
	"github.com/nynicg/cake/lib/ahoy"
	"github.com/nynicg/cake/lib/encrypt"
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/pool"
	"github.com/nynicg/cake/lib/socks5"
	"io"
	"net"
	"sync"
)

func startLocalSocksProxy(encryptor encrypt.Encryptor){
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

func handleCliConn(cliconn net.Conn ,pl *pool.TcpConnPool,encryptor encrypt.Encryptor){
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
	log.Debug("parse remote host -> " ,addr.Address())

	var remote net.Conn
	var bypass int
	if i ,ok := GetDomainCache(addr.Host());ok{
		bypass = i
	}else{
		bypass = Bypass(addr.Host())
		PutDomainCache(addr.Host() ,bypass)
	}
	var encryptorSelect encrypt.Encryptor
	switch bypass {
	case BypassDiscard:
		socks5.ProxyFailed(socks5.SocksRespHostUnreachable ,cliconn)
		return
	case BypassTrue:
		encryptorSelect ,_ = encrypt.GetStreamEncryptor(encrypt.EncryptTypePlain)
		remote ,e = net.Dial("tcp" ,addr.Address())
		if e != nil{
			log.Errorx("dail bypassed remote addr " ,e)
			socks5.ProxyFailed(socks5.SocksRespHostUnreachable ,cliconn)
			return
		}
	case BypassProxy:
		encryptorSelect = encryptor
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

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		upN ,e := ahoy.CopyWithEncryptor(remote ,cliconn ,encryptorSelect.Encrypt)
		//upN ,e := io.Copy(remote ,cliconn)
		if e != nil{
			log.Warn("copy client request to remote server error " ,e)
			return
		}
		log.Debug(fmt.Sprintf("%s %d bit ↑" ,remote.RemoteAddr() ,upN))
	}()
	go func() {
		defer wg.Done()
		downN ,e := ahoy.CopyWithEncryptor(cliconn ,remote ,encryptorSelect.Decrypt)
		//downN ,e := encryptorSelect.DecryptStream(cliconn ,remote)
		//downN ,e := io.Copy(cliconn ,remote)
		if e != nil{
			log.Warn("copy server response to client error " ,e)
			return
		}
		log.Debug(fmt.Sprintf("%s %d bit ↓" ,remote.RemoteAddr() ,downN))
	}()
	wg.Wait()
}


func handshakeRemote(remote net.Conn ,proxyhost string) error{
	if len(proxyhost) > 255 {
		return errors.New("host addr is too long(>255)")
	}
	req := ahoy.RemoteConnectRequest{
		Encryption: 	byte(encrypt.GetStreamEncryptorIndexByName(config.EncryptType)),
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
