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

	Key				string
}

var config Config

func init(){
	cfg := &Config{}
	flag.StringVar(&cfg.AccessKey ,"k" , "M5Rm2nmNyn1cg@ru" ,"remote proxy server access key")
	flag.StringVar(&cfg.LocalAddr ,"s" ,"0.0.0.0:1921" ,"local proxy listening address")
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
	startProxyServ(enmap)
}



func registryEncrypt(config Config) *encrypt.EncryptorMap{
	enmap := encrypt.NewEncryptorMap()

	enmap.Register(encrypt.EncryptTypePlain ,&encrypt.Plain{})

	gcm ,e := encrypt.NewAES128GCM(config.Key)
	if e != nil {
		panic(e)
	}
	enmap.Register(encrypt.EncryptTypeAES128GCM ,gcm)

	cc ,e := encrypt.NewChacha20Poly1305(config.Key)
	if e != nil {
		panic(e)
	}
	enmap.Register(encrypt.EncryptTypeCHACHA ,cc)
	return enmap
}
