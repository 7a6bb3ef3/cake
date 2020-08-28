// +build windows

package main

import (
	"log"

	"golang.org/x/sys/windows/registry"
)


func configure(addr string){
	key ,e := registry.OpenKey(registry.CURRENT_USER ,"Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings" ,registry.ALL_ACCESS)
	if e != nil{
		log.Println(e)
		return
	}
	defer key.Close()
	key.SetBinaryValue("ProxyEnable" ,[]byte{1})
	key.SetStringValue("ProxyServer" ,addr)
}