package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nynicg/cake/lib/cryptor"
	"github.com/nynicg/cake/lib/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LocalAddr		string
	LocalApi		string
	Uid       		string
	LogLevel		int
	MaxConn			int
	Help			bool

	Key				string
}

var config Config

func init(){
	cfg := &Config{}
	flag.StringVar(&cfg.Uid ,"user" , "M5Rm2nmNyn1cg@ru" ,"recommend use uuid")
	flag.StringVar(&cfg.LocalAddr ,"addr" ,"0.0.0.0:1921" ,"local proxy listening address")
	flag.StringVar(&cfg.LocalApi ,"api" ,"0.0.0.0:1922" ,"local api listening address")
	flag.IntVar(&cfg.LogLevel ,"l" ,int(zap.InfoLevel) ,"log level(from -1 to 5)")
	flag.IntVar(&cfg.MaxConn ,"n" ,2048 ,"the maximum number of proxy connections")
	flag.BoolVar(&cfg.Help ,"help" ,false ,"display help info")
	flag.StringVar(&cfg.Key ,"key" ,"BAby10nStAGec0atBAby10nStAGec0at" ,"cryption methods key")
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
	enmap := registryEncrypt(config)
	go runProxyServ(enmap)
	runApiServ()
}



func registryEncrypt(config Config) *cryptor.CryptorMap{
	enmap := cryptor.NewEncryptorMap()

	enmap.Register(cryptor.CryptTypePlain ,&cryptor.Plain{})

	gcm ,e := cryptor.NewAES128GCM(config.Key)
	if e != nil {
		panic(e)
	}
	enmap.Register(cryptor.CryptTypeAES128GCM ,gcm)

	cc ,e := cryptor.NewChacha20Poly1305(config.Key)
	if e != nil {
		panic(e)
	}
	enmap.Register(cryptor.CryptTypeCHACHA ,cc)
	return enmap
}
