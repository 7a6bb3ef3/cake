package main

import (
	"github.com/nynicg/cake/lib/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LocalAddr		string
	AccessKey       string
	LogLevel		zapcore.Level
}

var config Config

func init(){
	config = Config{
		LocalAddr: "0.0.0.0:1921",
		AccessKey: "M5Rm2nmNyn1cg@ru",
		LogLevel:  zap.InfoLevel,
	}
}

func main(){
	log.InitLog(config.LogLevel)
	startProxyServ()
}
