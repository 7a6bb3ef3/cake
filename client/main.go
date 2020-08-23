package main

import (
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
	LogLevel			zapcore.Level
}

var config Config

func init(){
	config = Config{
		RemoteExitAddr: "127.0.0.1:1921",
		RemoteAccessKey: "M5Rm2nmNyn1cg@ru",
		LocalSocksAddr: "127.0.0.1:1919",
		LocalHttpAddr: "127.0.0.1:1920",
		MaxLocalConnNum: 1024,

		LogLevel:		zap.InfoLevel,
	}
}

func main(){
	log.InitLog(config.LogLevel)
	loadPassrule()
	go startLocalHttpProxy()
	startLocalSocksProxy()
}

