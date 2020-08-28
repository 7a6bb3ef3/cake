package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/nynicg/cake/lib/log"
)

type Config struct {
	RemoteExitAddr  string
	Uid             string
	LocalSocksAddr  string
	LocalHttpAddr   string
	MaxLocalConnNum int
	LogLevel        string
	Help            bool
	EncryptType     string
	Key             string
	AutoConfigure	bool
	DisableGui		bool
}

var config Config

func init() {
	cfg := &Config{}
	flag.StringVar(&cfg.Uid, "user", "M5Rm2nmNyn1cg@ru", "recommend use uuid")
	flag.StringVar(&cfg.RemoteExitAddr, "proxy", "127.0.0.1:1921", "remote proxy server address")
	flag.StringVar(&cfg.LocalHttpAddr, "http", "127.0.0.1:1919", "local http proxy listening address")
	flag.StringVar(&cfg.LocalSocksAddr, "socks", "127.0.0.1:1920", "local SOKCKS5 proxy listening address")
	flag.StringVar(&cfg.LogLevel, "lvl", "info", "log level(from debug to fatal)")
	flag.IntVar(&cfg.MaxLocalConnNum, "n", 2048, "the maximum number of local connections")
	flag.BoolVar(&cfg.Help, "h", false, "display help info")
	flag.BoolVar(&cfg.DisableGui, "nonGui", false, "place an icon and menu in the notification area(windows ONLY)")
	flag.BoolVar(&cfg.AutoConfigure, "auto", false, "auto configure system proxy")
	flag.StringVar(&cfg.EncryptType, "cryptor", "aes128gcm", "supported encryption methods ,following is the supported list:\n {chacha|aes128gcm|plain}")
	flag.StringVar(&cfg.Key, "key", "BAby10nStAGec0atBAby10nStAGec0at", "cryption methods key(length must be 32)")
	flag.Parse()
	flag.Usage = usage
	config = *cfg
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:cakecli [OPTIONS]...")
	flag.PrintDefaults()
}

func main() {
	if config.Help {
		usage()
		return
	}
	log.InitLog(config.LogLevel)
	if !config.DisableGui && runtime.GOOS == "windows"{
		log.Info("Open as icon")
		RunAsIcon(func() {})
	}
	log.Info("Use cryptor ", config.EncryptType)
	loadPassrule()
	go startLocalHttpProxy()
	runLocalSocksProxy()
}
