package main

import "time"

type Config struct {
	LocalAddr		string
	AccessKey       string
}

var config Config

func init(){
	config = Config{
		LocalAddr: "0.0.0.0:1921",
		AccessKey: "M5Rm2nmNyn1cg@ru",
	}
}

func main(){
	go startProxyServ()
	time.Sleep(time.Minute * 20)
}
