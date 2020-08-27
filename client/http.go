package main

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/nynicg/cake/lib/log"
	"golang.org/x/net/proxy"
)

func startLocalHttpProxy() {
	prox := goproxy.NewProxyHttpServer()
	prox.Logger = log.GetAdaptLogger()
	prox.Tr = &http.Transport{
		// have to use Dial here,if using DialContext ,
		// goproxy can not proxy https request ,maybe there is some unknown problem in elazarl/goproxy lib ,
		Dial:            httpDial,
		IdleConnTimeout: time.Minute,
		Proxy:           nil,
	}

	httpprox := &http.Server{
		Addr:         config.LocalHttpAddr,
		Handler:      prox,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 20,
	}

	log.Info("HTTP listen on ", config.LocalHttpAddr)
	if err := httpprox.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func httpDial(network string, addr string) (net.Conn, error) {
	log.Debug("recv http proxy request to ", addr, network)
	conn, e := redirectToLocalSocks(addr)
	if e != nil {
		return nil, e
	}
	tcpcon, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, errors.New("assert failed")
	}
	tcpcon.SetKeepAlive(false)
	tcpcon.SetWriteBuffer(1 << 15)
	return conn, nil
}

func redirectToLocalSocks(addr string) (net.Conn, error) {
	socks, e := proxy.SOCKS5("tcp", config.LocalSocksAddr, nil, proxy.Direct)
	if e != nil {
		log.Panic(e)
	}
	return socks.Dial("tcp", addr)
}
