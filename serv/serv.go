package main

import (
	"fmt"
	"github.com/nynicg/cake/lib"
	"io"
	"net"
	"sync"
	"time"

	"github.com/nynicg/cake/lib/ahoy"
	"github.com/nynicg/cake/lib/log"
)

func startProxyServ() {
	ls ,e := net.Listen("tcp" ,config.LocalAddr)
	if e != nil {
		log.Panic(e)
	}
	pool := lib.NewTcpConnPool(config.MaxConn)
	log.Info("Listen on " ,config.LocalAddr)
	for {
		fromsocks := pool.GetLocalTcpConn()
		fromsocks ,e := ls.Accept()
		if e != nil{
			log.Errorx("accept new client conn " ,e)
			continue
		}
		go handleConn(fromsocks ,pool)
	}
}

func handleConn(fromsocks net.Conn ,pool *lib.TcpConnPool){
	log.Debug("handle conn from " ,fromsocks.RemoteAddr())
	defer func() {
		fromsocks.Close()
		pool.FreeLocalTcpConn(fromsocks)
	}()
	fromsocks.(*net.TCPConn).SetKeepAlive(false)
	addr ,e := ahoy.Handshake(config.AccessKey ,fromsocks)
	if e != nil{
		log.Errorx("handshake " ,e)
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
	outConn.SetWriteDeadline(time.Now().Add(time.Minute))
	// ready to mathx
	if e := ahoy.OnReady(fromsocks);e != nil{
		log.Errorx("done handshake " + addr ,e)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		upN ,e := io.Copy(outConn ,fromsocks)
		if e != nil{
			log.Warn("copy client request to remote server error " ,e)
			return
		}
		log.Debug(fmt.Sprintf("%s %d bit ↑" ,outConn.RemoteAddr() ,upN))

	}()

	go func() {
		defer wg.Done()
		downN ,e := io.Copy(fromsocks ,outConn)
		if e != nil{
			log.Warn("copy server response to client error " ,e)
			return
		}
		log.Debug(fmt.Sprintf("%s %d bit ↓" ,outConn.RemoteAddr() ,downN))
	}()
	wg.Wait()
}


