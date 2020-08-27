package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nynicg/cake/lib/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	RemoteExitAddr  string
	Uid             string
	LocalSocksAddr  string
	LocalHttpAddr   string
	MaxLocalConnNum int
	LogLevel        int
	Help            bool
	EncryptType     string
	Key             string
}

var config Config

func init() {
	cfg := &Config{}
	flag.StringVar(&cfg.Uid, "user", "M5Rm2nmNyn1cg@ru", "recommend use uuid")
	flag.StringVar(&cfg.RemoteExitAddr, "proxy", "127.0.0.1:1921", "remote proxy server address")
	flag.StringVar(&cfg.LocalHttpAddr, "http", "127.0.0.1:1919", "local http proxy listening address")
	flag.StringVar(&cfg.LocalSocksAddr, "socks", "127.0.0.1:1920", "local SOKCKS5 proxy listening address")
	flag.IntVar(&cfg.LogLevel, "lvl", int(zap.InfoLevel), "log level(from -1 to 5)")
	flag.IntVar(&cfg.MaxLocalConnNum, "n", 2048, "the maximum number of local connections")
	flag.BoolVar(&cfg.Help, "h", false, "display help info")
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
	log.InitLog(zapcore.Level(config.LogLevel))
	log.Info("Use cryptor ", config.EncryptType)
	loadPassrule()
	go startLocalHttpProxy()
	runLocalSocksProxy()
}
