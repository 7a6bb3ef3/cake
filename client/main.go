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
	ChachaKey			string
}

var config Config
var defEncryptor encrypt.Encryptor

func init(){
	cfg := &Config{}
	flag.StringVar(&cfg.RemoteAccessKey ,"k" , "M5Rm2nmNyn1cg@ru" ,"remote proxy server access key")
	flag.StringVar(&cfg.RemoteExitAddr ,"r" ,"127.0.0.1:1921" ,"remote proxy server address")
	flag.StringVar(&cfg.LocalHttpAddr ,"h" ,"127.0.0.1:1919" ,"local http proxy listening address")
	flag.StringVar(&cfg.LocalSocksAddr ,"s" ,"127.0.0.1:1920" ,"local SOKCKS5 proxy listening address")
	flag.IntVar(&cfg.LogLevel ,"l" ,int(zap.InfoLevel) ,"log level(from -1 to 5)")
	flag.IntVar(&cfg.MaxLocalConnNum ,"n" ,2048 ,"the maximum number of local connections")
	flag.BoolVar(&cfg.Help ,"help" ,false ,"display help info")
	flag.StringVar(&cfg.EncryptType ,"encrypt" ,"chacha" ,"supported encryption methods ,following is the supported list:\n {chacha|AES128CFB|PLAIN}")
	flag.StringVar(&cfg.AESKey ,"aesKey" ,"BAby10nStAGec0at" ,"key of AES cryption")
	flag.StringVar(&cfg.ChachaKey ,"chachaKey" ,"srMysu9kidEsuNeIcgnOCAkes1zanEki" ,"key of Chacha20poly1305")
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
	case "aes128cfb":
		cfb ,e := encrypt.NewAES128CFB(config.AESKey)
		if e != nil {
			panic(e)
		}
		defEncryptor = cfb
	case "chacha":
		cc ,e := encrypt.NewChacha20Poly1305(config.ChachaKey ,encrypt.DefaultChachaNonce ,encrypt.DefaultChachaAad)
		if e != nil {
			panic(e)
		}
		defEncryptor = cc
	case "plain":
		defEncryptor = &encrypt.Plain{}
	default:
		panic("unsupported encryption method")
	}
}


