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
)

var bufpool *pool.BufferPool
var proxyStat *ProxyStat

func init() {
	bufpool = pool.NewBufPool(16 * 1024)
	proxyStat = &ProxyStat{}
}

func runProxyServ() {
	ls, e := net.Listen("tcp4", globalConfig.ProxyConfig.LocalAddr)
	if e != nil {
		log.Panic(e)
	}
	pl := pool.NewTickets(globalConfig.ProxyConfig.MaxConn)
	log.Info("Listen on ", globalConfig.ProxyConfig.LocalAddr)
	for {
		pl.GetTicket()
		src, e := ls.Accept()
		if e != nil {
			log.Errorx("accept new client conn ", e)
			continue
		}
		go handleConn(src, pl)
	}
}

func handleConn(src net.Conn, pl *pool.Tickets) {
	log.Debug("handle conn from ", src.RemoteAddr())
	defer func() {
		src.Close()
		pl.FreeTicket()
	}()
	src.(*net.TCPConn).SetKeepAlive(false)
	info, e := handshake(src)
	if e != nil {
		log.Errorx("handshake ", e)
		return
	}
	crypt, e := cryptor.NewCryptor(info.cryptType, string(info.randomKey))
	if e != nil {
		log.Errorx("get stream encryptor ", e)
		return
	}
	dst, e := net.Dial("tcp", info.addr)
	if e != nil {
		log.Error("dial proxy addr ", info.addr, e)
		return
	}
	defer dst.Close()
	dst.(*net.TCPConn).SetKeepAlive(false)

	inboundEnv := &ahoy.CopyEnv{
		ReaderWithLength: false,
		WriterNeedLength: true,
		CryptFunc:        crypt.Encrypt,
		BufPool:          bufpool,
	}
	outboundEnv := &ahoy.CopyEnv{
		ReaderWithLength: true,
		WriterNeedLength: false,
		CryptFunc:        crypt.Decrypt,
		BufPool:          bufpool,
	}
	var (
		up   int
		down int
	)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		downN, e := ahoy.CopyConn(src, dst, inboundEnv)
		down = downN
		if e != nil {
			log.Info(info.addr, " server resp.", e)
			return
		}
	}()

	go func() {
		defer func() {
			wg.Done()
			dst.(*net.TCPConn).CloseWrite()
		}()
		upN, e := ahoy.CopyConn(dst, src, outboundEnv)
		up = upN
		if e != nil {
			log.Info(info.addr, " client request. ", e)
			return
		}
	}()
	wg.Wait()
	onFinish(up, down, info.addr)
}

type hsinfo struct {
	cryptType int
	randomKey []byte
	addr      string
}

// hmac 16byte
//
// use a customer protocol ,for experiment
// return encryption type ,proxy address and an error if there is
//  +-----+-----+-----+-----+-----+
//  |ENC  |CMD  |RDKey|LEN  |ADDR |
//  +-----+-----+-----+-----+-----+
//  |1    |1    |32   |1    |LEN  |
// if success ,server response(random 6 bit)
func handshake(src net.Conn) (hsinfo, error) {
	info := hsinfo{}
	buf := bufpool.Get()
	defer bufpool.Put(buf)
	// verify HMAC
	if _, e := io.ReadFull(src, buf[:16]); e != nil {
		return info, e
	}
	if DefUidManager().VerifyHMAC(buf[:16]) {
		src.Write([]byte{ahoy.HMACOK})
	} else {
		src.Write([]byte{ahoy.HMACInvalid})
		return info, errors.New("hmac auth failed")
	}

	// parse cmd
	if _, e := io.ReadFull(src, buf[:35]); e != nil {
		return info, e
	}
	p1 := buf[:35]
	pr := make([]byte, 35)
	cryptor.XorStream(pr, buf[:35], globalConfig.ProxyConfig.Key)
	addrLen := pr[34]
	info.cryptType = int(pr[0])
	info.randomKey = pr[2:34]
	if pr[1] != byte(ahoy.CommandConnect) {
		return info, errors.New("unsupport command")
	}
	if addrLen == 0 {
		return info, errors.New("empty proxy addr")
	}
	// read addr
	if _, e := io.ReadFull(src, buf[:addrLen]); e != nil {
		return info, e
	}
	cryptor.XorStream(buf, append(p1, buf[:addrLen]...), globalConfig.ProxyConfig.Key)
	info.addr = string(buf[35 : 35+addrLen])
	// TODO client need some specific msg, no one likes 114514
	_, e := src.Write([]byte{1, 1, 4, 5, 1, 4})
	if e != nil {
		return info, e
	}
	return info, nil
}

func onFinish(up, down int, addr string) {
	log.Debug(fmt.Sprintf("%s ,%d ↑ ,%d ↓ bytes", addr, up, down))
	proxyStat.Add(up, down)
}
