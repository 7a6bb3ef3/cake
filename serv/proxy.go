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
	bufpool = pool.NewBufPool(32 * 1024)
	proxyStat = &ProxyStat{}
}

func runProxyServ(enmap *cryptor.CryptorMap) {
	ls, e := net.Listen("tcp", config.LocalAddr)
	if e != nil {
		log.Panic(e)
	}
	pl := pool.NewTcpConnPool(config.MaxConn)
	log.Info("Listen on ", config.LocalAddr)
	for {
		fromsocks := pl.GetLocalTcpConn()
		fromsocks, e := ls.Accept()
		if e != nil {
			log.Errorx("accept new client conn ", e)
			continue
		}
		go handleConn(fromsocks, pl, enmap)
	}
}

func handleConn(fromsocks net.Conn, pl *pool.TcpConnPool, enmap *cryptor.CryptorMap) {
	log.Debug("handle conn from ", fromsocks.RemoteAddr())
	defer func() {
		fromsocks.Close()
		pl.FreeLocalTcpConn(fromsocks)
	}()
	fromsocks.(*net.TCPConn).SetKeepAlive(false)
	info, e := handshake(fromsocks)
	if e != nil {
		log.Errorx("handshake ", e)
		return
	}
	crypt, e := cryptor.NewCryptor(info.cryptType ,string(info.randomKey))
	if e != nil {
		log.Errorx("get stream encryptor ", e)
		return
	}
	outConn, e := net.Dial("tcp", info.addr)
	if e != nil {
		log.Error("dial proxy addr " ,info.addr, e)
		return
	}
	defer outConn.Close()
	outConn.(*net.TCPConn).SetKeepAlive(false)
	if e := onReady(fromsocks); e != nil {
		log.Error("done handshake " ,info.addr, e)
		return
	}

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
		downN, e := ahoy.CopyConn(fromsocks, outConn, inboundEnv)
		down = downN
		if e != nil {
			log.Info(info.addr, " server resp.", e)
			return
		}
	}()

	go func() {
		defer func() {
			wg.Done()
			outConn.(*net.TCPConn).CloseWrite()
		}()
		upN, e := ahoy.CopyConn(outConn, fromsocks, outboundEnv)
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
	cryptType		int
	randomKey		[]byte
	addr			string
}

// use a customer protocol ,for experiment
// return encryption type ,proxy address and an error if there is
//  +-----+-----+-----+-----+-----+
//  |ENC  |CMD  |RDKey|LEN  |ADDR |
//  +-----+-----+-----+-----+-----+
//  |1    |1    |32   |1    |LEN  |
// if success ,server response(random 6 bit)
func handshake(fromsocks net.Conn) (hsinfo , error) {
	info := hsinfo{}
	buf := bufpool.Get()
	defer bufpool.Put(buf)
	if _, e := io.ReadFull(fromsocks, buf[:35]); e != nil {
		return info, e
	}
	p1 := buf[:35]
	pr := make([]byte, 35)
	cryptor.XorStream(pr, buf[:35], config.Key)
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
	if _, e := io.ReadFull(fromsocks, buf[:addrLen]); e != nil {
		return info, e
	}
	cryptor.XorStream(buf, append(p1, buf[:addrLen]...), config.Key)
	info.addr = string(buf[35 : 35+addrLen])
	return info, nil
}

func onReady(w io.Writer) error {
	_, e := w.Write([]byte{1, 1, 4, 5, 1, 4})
	return e
}

func onFinish(up, down int, addr string) {
	log.Debug(fmt.Sprintf("%s ,%d ↑ ,%d ↓ bytes", addr, up, down))
	proxyStat.Add(up, down)
}
