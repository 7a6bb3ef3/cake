// +build windows

package main

import (
	"github.com/nynicg/cake/lib/log"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/sys/windows/registry"
)

var defaultRegVal = []byte{
	70 ,0 ,0 ,0 ,50 ,1 ,0 ,0 ,3 ,0 ,0 ,0 ,14 ,0 ,0 ,0 ,49 ,50 ,55 ,46 ,48 ,46 ,48 ,46 ,49 ,58 ,49 ,57 ,49 ,57 ,182 ,
	0 ,0 ,0 ,108 ,111 ,99 ,97 ,108 ,104 ,111 ,115 ,116 ,59 ,49 ,50 ,55 ,46 ,42 ,59 ,49 ,48 ,46 ,42 ,59 ,49 ,55 ,50 ,
	46 ,49 ,54 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,49 ,55 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,49 ,56 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,
	49 ,57 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,50 ,48 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,50 ,49 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,50 ,
	50 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,50 ,51 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,50 ,52 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,50 ,53 ,
	46 ,42 ,59 ,49 ,55 ,50 ,46 ,50 ,54 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,50 ,55 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,50 ,56 ,46 ,
	42 ,59 ,49 ,55 ,50 ,46 ,50 ,57 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,51 ,48 ,46 ,42 ,59 ,49 ,55 ,50 ,46 ,51 ,49 ,46 ,42 ,
	59 ,49 ,57 ,50 ,46 ,49 ,54 ,56 ,46 ,42 ,59 ,60 ,108 ,111 ,99 ,97 ,108 ,62 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,
	0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0 ,0}

func getDefaultConnSetting() ([]byte ,error){
	key ,_ ,e := registry.CreateKey(registry.CURRENT_USER ,"Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Connections" ,registry.ALL_ACCESS)
	if e != nil{
		return nil ,e
	}
	defer key.Close()
	raw ,_ ,e := key.GetBinaryValue("DefaultConnectionSettings")
	if e != nil{
		return raw ,e
	}
	if raw[4] == 255 {
		raw[4] = 0
	}else{
		raw[4] = raw[4] + 1
	}
	defaultRegVal[4] = raw[4]
	return defaultRegVal[:] ,nil
}

// configure modify the windows registry to enable the system proxy
func configure(){
	b ,e := getDefaultConnSetting()
	if e != nil{
		log.Error(e)
		return
	}
	b[8] = 3
	key ,_ ,e := registry.CreateKey(registry.CURRENT_USER ,"Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Connections" ,registry.ALL_ACCESS)
	if e != nil{
		log.Error(e)
		return
	}
	defer key.Close()
	key.SetBinaryValue("DefaultConnectionSettings" ,b)
	//refreshReg()
}


// unconfigure like configure ,modify the windows registry to disable the proxy
func unconfigure() {
	b ,e := getDefaultConnSetting()
	if e != nil{
		log.Error(e)
		return
	}
	b[8] = 1
	key ,_ ,e := registry.CreateKey(registry.CURRENT_USER ,"Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Connections" ,registry.ALL_ACCESS)
	if e != nil{
		log.Error(e)
		return
	}
	defer key.Close()
	key.SetBinaryValue("DefaultConnectionSettings" ,b)
	//refreshReg()
}

// TODO it will not take effect after updating the registry ,so we need to do something can refresh registry
func refreshReg(){
	open.Run("ms-settings:network-proxy")
}