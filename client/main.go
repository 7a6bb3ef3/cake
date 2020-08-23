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
	RemoteExitAddr 	string
	RemoteAccessKey		string
	LocalSocksAddr		string
	LocalHttpAddr		string
	MaxLocalConnNum		int
	LogLevel			int
	Help				bool
	EncryptType			string
	AESKey				string
	AESIv				string
}

var config Config

func init(){
	cfg := &Config{}
	flag.StringVar(&cfg.RemoteAccessKey ,"k" , "M5Rm2nmNyn1cg@ru" ,"remote proxy server access key")
	flag.StringVar(&cfg.RemoteExitAddr ,"r" ,"127.0.0.1:1921" ,"remote proxy server address")
	flag.StringVar(&cfg.LocalHttpAddr ,"h" ,"127.0.0.1:1919" ,"local http proxy listening address")
	flag.StringVar(&cfg.LocalSocksAddr ,"s" ,"127.0.0.1:1920" ,"local SOKCKS5 proxy listening address")
	flag.IntVar(&cfg.LogLevel ,"l" ,int(zap.InfoLevel) ,"log level(from -1 to 5)")
	flag.IntVar(&cfg.MaxLocalConnNum ,"n" ,2048 ,"the maximum number of local connections")
	flag.BoolVar(&cfg.Help ,"help" ,false ,"display help info")
	flag.StringVar(&cfg.EncryptType ,"encrypt" ,"AES128CBC" ,"supported encryption methods ,following is the supported list:\n {AES128CBC|AES128CFB|PLAIN}")
	flag.StringVar(&cfg.AESKey ,"aesKey" ,"BAby10nStAGec0at" ,"key of AES cryption")
	flag.StringVar(&cfg.AESIv ,"aesIv" ,"j0ker_nE1_diyusi" ,"vi of AES_CBC cryption")
	flag.Parse()
	flag.Usage = usage
	config = *cfg
}

func usage(){
	fmt.Fprintln(os.Stderr ,"Usage:cakecli [OPTIONS]...")
	flag.PrintDefaults()
}

func main(){
	if config.Help {
		usage()
		return
	}
	log.InitLog(zapcore.Level(config.LogLevel))
	loadPassrule()
	go startLocalHttpProxy()
	en ,e := encrypt.NewStreamEncryptorByName(config.EncryptType ,config.AESKey ,config.AESIv)
	if e != nil{
		log.Panic(e)
	}
	log.Info("Use encryption " ,config.EncryptType)
	startLocalSocksProxy(en)
}


