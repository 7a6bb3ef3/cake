package main

import (
	"github.com/nynicg/cake/lib/cryptor"
	"sync"
	"time"
)

var defUidMg *UIDManager

func init(){
	defUidMg = &UIDManager{
		uidMap: make(map[string]UIDInfo),
		hmacCache: make(map[string]int64),
		hmactmp: make(map[string]int64),
	}
	tk := time.NewTicker(time.Second * 8)
	go func() {
		for range tk.C{
			defUidMg.refreshCache()
		}
	}()
}

func loadUidsFromCfg(){
	for _ ,v := range globalConfig.ProxyConfig.Uids{
		DefUidManager().RegisterUid(v ,UIDInfo{
			CreateTime: time.Now().Unix(),
			Addr:       "from configure",
		})
	}
	defUidMg.refreshCache()
}

func DefUidManager() *UIDManager{
	return defUidMg
}

type UIDInfo struct {
	CreateTime		int64
	Addr			string
}

type UIDManager struct {
	uidMap 		map[string]UIDInfo
	m			sync.Mutex

	// [hmac hex]create time
	hmacCache	map[string]int64
	hmactmp		map[string]int64
	cm			sync.Mutex
}

func (u *UIDManager)RegisterUid(uid string ,info UIDInfo){
	u.m.Lock()
	u.uidMap[uid] = info
	u.m.Unlock()
}

func (u *UIDManager)UnregisterUid(uid string){
	u.m.Lock()
	delete(u.uidMap ,uid)
	u.m.Unlock()
}

func (u *UIDManager)VerifyHMAC(hmac []byte) bool{
	u.cm.Lock()
	 _ ,ok := u.hmacCache[string(hmac)]
	u.cm.Unlock()
	return ok
}

func (u *UIDManager)refreshCache(){
	u.m.Lock()
	ulist := make([]string ,len(u.uidMap))
	i := 0
	for k := range u.uidMap{
		ulist[i] = k
		i++
	}
	u.m.Unlock()

	for _ ,v := range ulist{
		for _ ,validHmac := range cryptor.HMACAllTime(v){
			// use hex.encodetostring here cause more mem allocation
			u.hmactmp[string(validHmac)] = time.Now().Unix()
		}
	}

	u.cm.Lock()
	u.hmacCache = u.hmactmp
	u.cm.Unlock()

	u.hmactmp = make(map[string]int64)
}
