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

func init() {
	bufpool = pool.NewBufPool(32 * 1024)
}

func runLocalSocksProxy() {
	ls, e := net.Listen("tcp", config.LocalSocksAddr)
	if e != nil {
		log.Panic(e)
	}
	log.Info("Socks5 listen on ", config.LocalSocksAddr)
	for {
		src, e := ls.Accept()
		if e != nil {
			log.Errorx("accept new client conn ", e)
			continue
		}
		go handleCliConn(src)
	}
}

func handleCliConn(src net.Conn) {
	defer src.Close()
	src.(*net.TCPConn).SetKeepAlive(false)
	if e := socks5.Handshake(src); e != nil {
		log.Errorx("handshake with "+src.RemoteAddr().String(), e)
		return
	}
	addr, e := socks5.ParseCMD(src)
	if e != nil {
		log.Errorx("parse client cmd and addr", e)
		return
	}

	bypass := Bypass(addr.Host())
	var (
		dst           net.Conn
		cryptorSelect cryptor.Cryptor
	)
	switch bypass {
	case BypassDiscard:
		socks5.ProxyFailed(socks5.SocksRespHostUnreachable, src)
		return
	case BypassTrue:
		cryptorSelect = cryptor.GetTypePlain()
		dst, e = net.Dial("tcp", addr.Address())
		if e != nil {
			log.Errorx("dail bypassed remote addr ", e)
			socks5.ProxyFailed(socks5.SocksRespHostUnreachable, src)
			return
		}
	case BypassProxy:
		i, err := cryptor.GetIndexByName(config.EncryptType)
		if err != nil {
			log.Errorx("dail bypassed remote addr ", e)
			socks5.ProxyFailed(socks5.SocksRespHostUnreachable, src)
			return
		}
		cp, rdk, e := cryptor.NewCryptorX(i)
		if e != nil {
			log.Errorx("newCryptorX ", e)
			socks5.ProxyFailed(socks5.SocksRespServErr, src)
			return
		}
		cryptorSelect = cp
		dst, e = net.Dial("tcp", config.RemoteExitAddr)
		if e != nil {
			log.Errorx("dail remote exit ", e)
			socks5.ProxyFailed(socks5.SocksRespServErr, src)
			return
		}
		if e := handshakeRemote(dst, addr.Address(), rdk); e != nil {
			log.Errorx("handshake with remote failed ", e)
			socks5.ProxyFailed(socks5.SocksRespServErr, src)
			return
		}
	default:
		socks5.ProxyFailed(socks5.SocksRespServErr, src)
		return
	}
	defer dst.Close()
	dst.(*net.TCPConn).SetKeepAlive(false)
	if e := socks5.ProxyOK(src); e != nil {
		log.Errorx("local socks5 sent OK resp ", e)
		return
	}

	outboundEnv := &ahoy.CopyEnv{
		ReaderWithLength: false,
		WriterNeedLength: true,
		CryptFunc:        cryptorSelect.Encrypt,
		BufPool:          bufpool,
		Bypass:           bypass == BypassTrue,
	}
	inboundEnv := &ahoy.CopyEnv{
		ReaderWithLength: true,
		WriterNeedLength: false,
		CryptFunc:        cryptorSelect.Decrypt,
		BufPool:          bufpool,
		Bypass:           bypass == BypassTrue,
	}

	var (
		up   int
		down int
	)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer func() {
			wg.Done()
			dst.(*net.TCPConn).CloseWrite()
		}()
		upn, e := ahoy.CopyConn(dst, src, outboundEnv)
		up = upn
		if e != nil {
			log.Info(addr.Address(), " src request.", e)
			return
		}
	}()
	go func() {
		defer wg.Done()
		downn, e := ahoy.CopyConn(src, dst, inboundEnv)
		down = downn
		if e != nil {
			log.Info(addr.Address(), " dst resp.", e)
			return
		}
	}()
	wg.Wait()
	onFinish(up, down, addr.Host())
}

func handshakeRemote(remote net.Conn, proxyhost string, rdk string) error {
	if len(proxyhost) > 255 {
		return errors.New("host addr is too long(>255)")
	}
	index, e := cryptor.GetIndexByName(config.EncryptType)
	if e != nil {
		return e
	}
	req := ahoy.RemoteConnectRequest{
		Encryption: byte(index),
		Command:    ahoy.CommandConnect,
		RandomKey:  []byte(rdk),
		AddrLength: byte(len(proxyhost)),
		Addr:       []byte(proxyhost),
	}
	bts, e := req.Bytes()
	if e != nil {
		return e
	}
	cryptor.XorStream(bts, bts, config.Key)
	if _, e := remote.Write(bts); e != nil {
		return e
	}
	// read ok resp {1,1,4,5,1,4}
	buf := make([]byte, 6)
	if _, e := io.ReadFull(remote, buf); e != nil {
		return e
	}
	return nil
}

func onFinish(up, down int, addr string) {
	log.Debug(fmt.Sprintf("%s ,%d ↑ ,%d ↓ bytes", addr, up, down))
}
