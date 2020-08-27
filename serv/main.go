package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"os"

	"github.com/nynicg/cake/lib/cryptor"
	"github.com/nynicg/cake/lib/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LocalAddr string
	Uid       string
	LogLevel  int
	MaxConn   int
	Help      bool
	Key       string
	// api server
	EnableAPI    bool
	LocalApiAddr string
	BAUserName   string
	BAPassword   string
}

var config Config

func init() {
	cfg := &Config{}
	flag.StringVar(&cfg.Uid, "user", "M5Rm2nmNyn1cg@ru", "recommend use uuid")
	flag.StringVar(&cfg.LocalAddr, "addr", "0.0.0.0:1921", "local proxy listening address")
	flag.IntVar(&cfg.LogLevel, "l", int(zap.InfoLevel), "log level(from -1 to 5)")
	flag.IntVar(&cfg.MaxConn, "n", 2048, "the maximum number of proxy connections")
	flag.BoolVar(&cfg.Help, "h", false, "display help info")
	flag.StringVar(&cfg.Key, "key", "BAby10nStAGec0atBAby10nStAGec0at", "cryption methods key")

	flag.StringVar(&cfg.LocalApiAddr, "apiAddr", "0.0.0.0:1922", "local api listening address")
	flag.BoolVar(&cfg.EnableAPI, "api", false, "enable api service")
	flag.StringVar(&cfg.BAUserName, "apiUser", uuid.New().String(), "base auth user name(random initial value)")
	flag.StringVar(&cfg.BAPassword, "apiPwd", uuid.New().String(), "base auth password(random initial value)")
	flag.Parse()
	flag.Usage = usage
	config = *cfg
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:cakeserv [OPTIONS]...")
	flag.PrintDefaults()
}

func main() {
	if config.Help {
		usage()
		return
	}
	log.InitLog(zapcore.Level(config.LogLevel))
	enmap := cryptor.RegistryAllCrypto(config.Key)
	go runApiServ()
	runProxyServ(enmap)
}