package main

import (
	"errors"
	"fmt"
	"github.com/nynicg/cake/lib/ahoy"
	"github.com/nynicg/cake/lib/log"
	"io"
	"net"
)

func startProxyServ() {
	ls ,e := net.Listen("tcp" ,":1921")
	if e != nil {
		log.Panic(e)
	}
	log.Info("Listen on " ,config.LocalAddr)
	for {
		fromsocks ,e := ls.Accept()
		if e != nil{
			log.Panic(e)
		}
		go handleConn(fromsocks)
	}
}

func handleConn(fromsocks net.Conn){
	log.Debug("handle conn from " ,fromsocks.RemoteAddr())
	defer fromsocks.Close()
	fromsocks.(*net.TCPConn).SetKeepAlive(false)
	addr ,e := handshake(config.AccessKey ,fromsocks)
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
	// ready to transport
	if e := onReady(fromsocks);e != nil{
		log.Errorx("done handshake " + addr ,e)
		return
	}

	log.Debug("ready to transport " ,fromsocks.RemoteAddr())
	go func() {
		upN ,e := io.Copy(outConn ,fromsocks)
		if e != nil{
			log.Warn("copy client request to remote server error " ,e)
			return
		}
		log.Debug(fmt.Sprintf("%s %d bit ↑" ,outConn.RemoteAddr() ,upN))
	}()

	downN ,e := io.Copy(fromsocks ,outConn)
	if e != nil{
		log.Warn("copy server response to client error " ,e)
		return
	}
	log.Debug(fmt.Sprintf("%s %d bit ↓" ,outConn.RemoteAddr() ,downN))
}

func handshake(ackey string ,fromsocks net.Conn) (string ,error){
	buf := make([]byte ,255)
	if _ ,e := io.ReadFull(fromsocks ,buf[:19]);e != nil{
		return "" ,e
	}
	addrLen := buf[18]
	if buf[1] != byte(ahoy.CommandConnect) {
		return "" ,errors.New("unsupport command")
	}else if string(buf[2:18]) != ackey {
		return "" ,errors.New("access refused")
	}else if addrLen == 0{
		return "" ,errors.New("empty proxy addr")
	}
	// read addr
	if _ ,e := io.ReadFull(fromsocks ,buf[:addrLen]);e != nil{
		return "" ,e
	}
	return string(buf[:addrLen]) ,nil
}

func onReady(w io.Writer) error{
	_ ,e := w.Write([]byte{1 ,1 ,4 ,5 ,1 ,4})
	return e
}
