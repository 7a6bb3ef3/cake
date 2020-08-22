package main

import (
	"time"
)

type Config struct {
	RemoteExitAddr 	string

	RemoteAccessKey		string
	LocalSocksAddr		string
	LocalHttpAddr		string
	MaxLocalConnNum		int
}

var config Config

func init(){
	config = Config{
		RemoteExitAddr: "127.0.0.1:1921",
		RemoteAccessKey: "M5Rm2nmNyn1cg@ru",
		LocalSocksAddr: "127.0.0.1:1919",
		LocalHttpAddr: "127.0.0.1:1920",
		MaxLocalConnNum: 1024,
	}
}

func main(){
	loadApnic()
	go startLocalHttpProxy()
	go startLocalSocksProxy()
	// TODO use stop chan sign to exit program
	time.Sleep(time.Minute * 20)
}

