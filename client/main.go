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
	RemoteExitAddr 	string
	RemoteAccessKey		string
	LocalSocksAddr		string
	LocalHttpAddr		string
	MaxLocalConnNum		int
	LogLevel			int
	Help				bool
	EncryptType			string
	Key					string
}

var config Config
var defEncryptor cryptor.Cryptor

func init(){
	cfg := &Config{}
	flag.StringVar(&cfg.RemoteAccessKey ,"k" , "M5Rm2nmNyn1cg@ru" ,"remote proxy server access key")
	flag.StringVar(&cfg.RemoteExitAddr ,"r" ,"127.0.0.1:1921" ,"remote proxy server address")
	flag.StringVar(&cfg.LocalHttpAddr ,"h" ,"127.0.0.1:1919" ,"local http proxy listening address")
	flag.StringVar(&cfg.LocalSocksAddr ,"s" ,"127.0.0.1:1920" ,"local SOKCKS5 proxy listening address")
	flag.IntVar(&cfg.LogLevel ,"l" ,int(zap.InfoLevel) ,"log level(from -1 to 5)")
	flag.IntVar(&cfg.MaxLocalConnNum ,"n" ,2048 ,"the maximum number of local connections")
	flag.BoolVar(&cfg.Help ,"help" ,false ,"display help info")
	flag.StringVar(&cfg.EncryptType ,"cryptor" ,"chacha" ,"supported encryption methods ,following is the supported list:\n {chacha|AES128CFB|PLAIN}")
	flag.StringVar(&cfg.Key ,"key" ,"BAby10nStAGec0atBAby10nStAGec0at" ,"cryption methods key")
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
	log.Info("Use encryption " ,config.EncryptType)
	loadPassrule()
	setEncryptor(config)
	go startLocalHttpProxy()
	startLocalSocksProxy(defEncryptor)
}


func setEncryptor(config Config){
	switch config.EncryptType {
	case cryptor.NameAES128GCM:
		cfb ,e := cryptor.NewAES128GCM(config.Key)
		if e != nil {
			panic(e)
		}
		defEncryptor = cfb
	case cryptor.NameCHACHA:
		cc ,e := cryptor.NewChacha20Poly1305(config.Key)
		if e != nil {
			panic(e)
		}
		defEncryptor = cc
	case cryptor.NamePlain:
		defEncryptor = &cryptor.Plain{}
	default:
		panic("unsupported encryption method")
	}
}


