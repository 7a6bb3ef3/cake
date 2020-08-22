package main

import "time"

type Config struct {
	RemoteSocksAddr 	string

	LocalSocksAddr		string
	LocalHttpAddr		string
	MaxLocalConnNum		int
}

var config Config

func init(){
	config = Config{
		LocalSocksAddr: "127.0.0.1:1919",
		LocalHttpAddr: "127.0.0.1:1920",
		MaxLocalConnNum: 1024,
	}
}

func main(){
	loadApnic()
	go startLocalHttpProxy()
	go startLocalSocksProxy()
	time.Sleep(time.Minute * 20)
}

