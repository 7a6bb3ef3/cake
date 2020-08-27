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

func init(){
	bufpool = pool.NewBufPool(32 * 1024)
	proxyStat = &ProxyStat{}
}

func runProxyServ(enmap *cryptor.CryptorMap) {
	ls ,e := net.Listen("tcp" ,config.LocalAddr)
	if e != nil {
		log.Panic(e)
	}
	pl := pool.NewTcpConnPool(config.MaxConn)
	log.Info("Listen on " ,config.LocalAddr)
	for {
		fromsocks := pl.GetLocalTcpConn()
		fromsocks ,e := ls.Accept()
		if e != nil{
			log.Errorx("accept new client conn " ,e)
			continue
		}
		go handleConn(fromsocks ,pl ,enmap)
	}
}

func handleConn(fromsocks net.Conn ,pl *pool.TcpConnPool ,enmap *cryptor.CryptorMap){
	log.Debug("handle conn from " ,fromsocks.RemoteAddr())
	defer func() {
		fromsocks.Close()
		pl.FreeLocalTcpConn(fromsocks)
	}()
	fromsocks.(*net.TCPConn).SetKeepAlive(false)
	cryptType ,addr ,e := handshake(fromsocks)
	if e != nil{
		log.Errorx("handshake " ,e)
		return
	}
	log.Debug("got enctype and addr " ,cryptType ," " ,addr)
	crypt ,e := enmap.Get(cryptType)
	if e != nil{
		log.Errorx("get stream encryptor " ,e)
		return
	}
	outConn ,e := net.Dial("tcp" ,addr)
	if e != nil{
		log.Errorx("dial proxy addr " + addr ,e)
		return
	}
	defer outConn.Close()
	outConn.(*net.TCPConn).SetKeepAlive(false)
	if e := onReady(fromsocks);e != nil{
		log.Errorx("done handshake " + addr ,e)
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
		up int
		down int
	)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		downN ,e := ahoy.CopyConn(fromsocks ,outConn ,inboundEnv)
		down = downN
		if e != nil{
			log.Info(addr ," server resp." ,e)
			return
		}
	}()

	go func() {
		defer func() {
			wg.Done()
			outConn.(*net.TCPConn).CloseWrite()
		}()
		upN ,e := ahoy.CopyConn(outConn ,fromsocks ,outboundEnv)
		up = upN
		if e != nil{
			log.Info(addr ," client request. " ,e)
			return
		}
	}()
	wg.Wait()
	onFinish(up ,down ,addr)
}

// use a customer protocol ,for experiment
// return encryption type ,proxy address and an error if there is
func handshake(fromsocks net.Conn) (int ,string ,error){
	buf := bufpool.Get()
	defer bufpool.Put(buf)
	if _ ,e := io.ReadFull(fromsocks ,buf[:19]);e != nil{
		return 0 ,"" ,e
	}
	p1 := buf[:19]
	pr := make([]byte ,19)
	cryptor.XorStream(pr ,buf[:19] ,config.Key)
	addrLen := pr[18]
	enctype := pr[0]
	if pr[1] != byte(ahoy.CommandConnect) {
		return 0 ,"" ,errors.New("unsupport command")
	}
	//if !AuthHMAC(pr[2:18]) {
	//	return 0 ,"" ,errors.New("incorrect uid or command")
	//}
	if addrLen == 0{
		return 0 ,"" ,errors.New("empty proxy addr")
	}
	// read addr
	if _ ,e := io.ReadFull(fromsocks ,buf[:addrLen]);e != nil{
		return 0 ,"" ,e
	}
	cryptor.XorStream(buf ,append(p1 ,buf[:addrLen]...) ,config.Key)
	return int(enctype) ,string(buf[19:19+addrLen]) ,nil
}

func onReady(w io.Writer) error{
	_ ,e := w.Write([]byte{1 ,1 ,4 ,5 ,1 ,4})
	return e
}

func onFinish(up ,down int ,addr string){
	log.Debug(fmt.Sprintf("%s ,%d ↑ ,%d ↓ bytes" ,addr ,up ,down))
	proxyStat.Add(up ,down)
}


