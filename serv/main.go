package main

import (
	"flag"
	"fmt"
	"github.com/nynicg/cake/lib/encrypt"
	"github.com/nynicg/cake/lib/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	LocalAddr		string
	AccessKey       string
	LogLevel		int
	MaxConn			int
	Help			bool

	AESKey			string
	AESVi			string
}

var config Config

func init(){
	cfg := &Config{}
	flag.StringVar(&cfg.AccessKey ,"k" , "M5Rm2nmNyn1cg@ru" ,"remote proxy server access key")
	flag.StringVar(&cfg.LocalAddr ,"s" ,"0.0.0.0:1921" ,"local proxy listening address")
	flag.IntVar(&cfg.LogLevel ,"l" ,int(zap.InfoLevel) ,"log level(from -1 to 5)")
	flag.IntVar(&cfg.MaxConn ,"n" ,2048 ,"the maximum number of proxy connections")
	flag.BoolVar(&cfg.Help ,"help" ,false ,"display help info")
	flag.StringVar(&cfg.AESKey ,"aesKey" ,"BAby10nStAGec0at" ,"key of AES cryption")
	flag.StringVar(&cfg.AESVi ,"aesVi" ,"j0ker_nE1_diyusi" ,"vi of AES_CBC cryption")
	flag.Parse()
	flag.Usage = usage
	config = *cfg
}

func usage(){
	fmt.Fprintln(os.Stderr ,"Usage:cakeserv [OPTIONS]...")
	flag.PrintDefaults()
}

func main(){
	if config.Help {
		usage()
		return
	}
	log.InitLog(zapcore.Level(config.LogLevel))
	if e := initEncryptors();e != nil{
		log.Panic(e)
	}
	startProxyServ()
}

func initEncryptors() error{
	if e := encrypt.SetDefaultAES128CBC(config.AESKey ,config.AESVi);e != nil{
		return e
	}
	return nil
}
