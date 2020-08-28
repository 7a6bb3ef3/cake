// +build windows

package main

import (
	"log"

	"golang.org/x/sys/windows/registry"
)


// configure modify the windows registry to enable the system proxy
func configure(){
	key ,e := registry.OpenKey(registry.CURRENT_USER ,"Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings" ,registry.ALL_ACCESS)
	if e != nil{
		log.Println(e)
		return
	}
	defer key.Close()
	key.SetBinaryValue("ProxyEnable" ,[]byte{1})
	key.SetStringValue("ProxyServer" ,config.LocalHttpAddr)
	refreshReg()
}


// unconfigure like configure ,modify the windows registry to disable the proxy
func unconfigure() {
	key ,e := registry.OpenKey(registry.CURRENT_USER ,"Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings" ,registry.ALL_ACCESS)
	if e != nil{
		log.Println(e)
		return
	}
	defer key.Close()
	key.SetBinaryValue("ProxyEnable" ,[]byte{0})
	refreshReg()
}

// TODO it will not take effect after updating the registry ,so we need to do something can refresh registry
func refreshReg(){

}