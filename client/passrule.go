package main

import (
	"bufio"
	"github.com/nynicg/cake/lib/log"
	"net/http"
	"os"
	"strings"
)

const ApnicUrl = "http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"
const ipFile = "apnic-latest.txt"
const gfwlist = "yjsp.txt"

func loadApnic(){
	info ,e := os.Stat(ipFile)
	// TODO check sum
	if e == nil && !info.IsDir(){
		log.Info("Find " ,ipFile)
		return
	}

	log.Info("Loading ip list from " ,ApnicUrl)
	log.Info("This process may take several minutes")
	f ,e := os.OpenFile(ipFile ,os.O_CREATE | os.O_RDWR | os.O_APPEND ,0755)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	resp ,e := http.Get(ApnicUrl)
	if e != nil{
		panic(e)
	}
	rd := bufio.NewReader(resp.Body)
	for {
		line ,e := rd.ReadString('\n')
		if e != nil{
			break
		}
		if strings.Contains(line ,"ipv4") && strings.Contains(line ,"CN") {
			f.WriteString(line)
		}
	}
}

func Bypass(host string) bool {
	return false
}
