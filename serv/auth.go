package main

import (
	"encoding/hex"
	"fmt"
	"github.com/nynicg/cake/lib/ahoy"
	"github.com/nynicg/cake/lib/cryptor"
	"sync"
)

// map[hmac_hex]uid
var usrMap map[string]string
var mux sync.Mutex

func init(){
	usrMap = make(map[string]string)

	// test
	if e := RegisterUidCmd(ahoy.CommandConnect ,"M5Rm2nmNyn1cg@ru");e != nil{
		panic(e)
	}
}

// check whether the user has sent allowed commands
func AuthHMAC(hmcbit []byte) bool{
	mux.Lock()
	defer mux.Unlock()
	_ ,ok := usrMap[hex.EncodeToString(hmcbit)]
	return ok
}

func RegisterUidCmd(cmd ahoy.Command ,uid string) error{
	out, e := cryptor.NewHMAC(uid).SumAhoyHandshake(byte(cmd) ,uid ,ahoy.HMACLength)
	if e != nil{
		return fmt.Errorf("registerAuth.%w" ,e)
	}
	mux.Lock()
	defer mux.Unlock()
	usrMap[hex.EncodeToString(out)] = uid
	return nil
}
