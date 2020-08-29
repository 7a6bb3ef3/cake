package main

import (
	"flag"
	"fmt"
	"github.com/nynicg/cake/lib/log"
	"os"
)

func parse() bool{
	help := flag.Bool("h", false, "display help info")

	cfg := &ServConfig{}
	flag.StringVar(&cfg.ProxyConfig.LocalAddr, "addr", "0.0.0.0:1921", "local proxy listening address")
	flag.StringVar(&cfg.ProxyConfig.LogLevel, "lvl", "info", "log level(from debug to fatal)")
	flag.IntVar(&cfg.ProxyConfig.MaxConn, "n", 2048, "the maximum number of proxy connections")
	flag.StringVar(&cfg.ProxyConfig.Key, "key", "BAby10nStAGec0atBAby10nStAGec0at", "cryption methods key")
	flag.StringVar(&cfg.ApiConfig.LocalApiAddr, "apiAddr", "0.0.0.0:1922", "local api listening address")
	flag.BoolVar(&cfg.ApiConfig.EnableApi, "api", false, "enable api service")
	flag.StringVar(&cfg.ApiConfig.BasicAuthUser, "apiUser", "", "base auth user name(random initial value)")
	flag.StringVar(&cfg.ApiConfig.BasicAuthPassword, "apiPwd", "", "base auth password(random initial value)")
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:cakeserv [OPTIONS]...")
		flag.PrintDefaults()
	}
	override(globalConfig ,cfg)
	return *help
}

func main() {
	if parse() {
		flag.Usage()
		return
	}
	log.InitLog(globalConfig.ProxyConfig.LogLevel)
	loadUidsFromCfg()
	go runApiServ()
	runProxyServ()
}
