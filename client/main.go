package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/nynicg/cake/lib/log"
)

var config *ProxyCfg

// return help
func parse() bool {
	help := flag.Bool("h", false, "display help info")

	cfg := &ProxyCfg{}
	flag.StringVar(&cfg.LocalHttpAddr, "http", "", "local http proxy listening address")
	flag.StringVar(&cfg.ServerPerfer, "server", "", "server addr")
	flag.StringVar(&cfg.LocalSocksAddr, "socks", "", "local SOKCKS5 proxy listening address")
	flag.StringVar(&cfg.LogLevel, "lvl", "", "log level(from debug to fatal)")
	flag.BoolVar(&cfg.DisableGui, "nonGui", false, "place an icon and menu in the notification area(windows ONLY)")
	flag.BoolVar(&cfg.EnforceProxy, "enforce", false, "proxy for every conns")
	flag.StringVar(&cfg.EncryptType, "cryptor", "", "supported encryption methods ,following is the supported list:\n {chacha|aes128gcm|plain}")
	flag.StringVar(&cfg.Key, "key", "", "cryption methods key ,length must be 32 (16 byte)")
	flag.StringVar(&cfg.Uid, "uid", "", "user uuid ,length must be 32 (16 byte)")

	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:cakecli [OPTIONS]...")
		flag.PrintDefaults()
	}
	overrideByCmd(globCfg, cfg)
	config = &globCfg.ProxyCfg
	return *help
}

func main() {
	if parse() {
		flag.Usage()
		return
	}
	log.InitLog(config.LogLevel)
	if !config.DisableGui && runtime.GOOS == "windows" {
		log.Info("Open as icon")
		RunAsIcon()
	}
	log.Info("Use cryptor ", config.EncryptType)
	go startLocalHttpProxy()
	runLocalSocksProxy()
}
