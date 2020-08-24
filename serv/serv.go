package main

import (
	"fmt"
	"github.com/nynicg/cake/lib/ahoy"
	"github.com/nynicg/cake/lib/encrypt"
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/pool"
	"net"
	"sync"
)

var bufpool *pool.BufferPool

func init(){
	bufpool = pool.NewBufPool(64 * 1024)
}

func startProxyServ(enmap *encrypt.EncryptorMap) {
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

func handleConn(fromsocks net.Conn ,pl *pool.TcpConnPool ,enmap *encrypt.EncryptorMap){
	log.Debug("handle conn from " ,fromsocks.RemoteAddr())
	defer func() {
		fromsocks.Close()
		pl.FreeLocalTcpConn(fromsocks)
	}()
	fromsocks.(*net.TCPConn).SetKeepAlive(false)
	encryptType ,addr ,e := ahoy.Handshake(config.AccessKey ,fromsocks)
	if e != nil{
		log.Errorx("handshake " ,e)
		return
	}
	log.Debug("got encrypt type " ,encryptType)
	//encryptor ,e := encrypt.GetStreamEncryptor(encryptType)
	encryptor ,e := enmap.Get(encryptType)
	if e != nil{
		log.Errorx("get stream encryptor " ,e)
		return
	}
	log.Debug("get proxy addr " ,addr ," from remote " ,fromsocks.RemoteAddr())
	outConn ,e := net.Dial("tcp" ,addr)
	if e != nil{
		log.Errorx("dial proxy addr " + addr ,e)
		return
	}
	defer outConn.Close()
	outConn.(*net.TCPConn).SetKeepAlive(false)
	// ready to mathx
	if e := ahoy.OnReady(fromsocks);e != nil{
		log.Errorx("done handshake " + addr ,e)
		return
	}

	bufa := bufpool.Get()
	bufb := bufpool.Get()
	defer func() {
		bufpool.Put(bufa)
		bufpool.Put(bufb)
	}()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		upN ,e := ahoy.CopyWithCryptFunc(outConn ,fromsocks ,encryptor.Decrypt ,bufa)
		if e != nil{
			log.Warn("proxy request -> server." ,e)
			return
		}
		log.Debug(fmt.Sprintf("%s %d bit ↑" ,outConn.RemoteAddr() ,upN))

	}()

	go func() {
		defer wg.Done()
		downN ,e := ahoy.CopyWithCryptFunc(fromsocks ,outConn ,encryptor.Encrypt ,bufb)
		if e != nil{
			log.Warn("server resp -> client. " ,e)
			return
		}
		log.Debug(fmt.Sprintf("%s %d bit ↓" ,outConn.RemoteAddr() ,downN))
	}()
	wg.Wait()
}


