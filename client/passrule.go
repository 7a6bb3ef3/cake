package main

import (
	"bufio"
	"github.com/BurntSushi/toml"
	"net/http"
	"os"
	"strings"

	_ "github.com/BurntSushi/toml"
	"github.com/nynicg/cake/lib/log"
)

const ApnicUrl = "http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"

const (
	domainFile = "domain.toml"
	ipFile = "apnic-latest.txt"
)

var cniplist = make(map[string]struct{})
var domainlist domain

func loadApnic(){
	defer func() {
		go func() {
			loadCNIPList()
			domainlist = loadDomainList()
		}()
	}()
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
	f.WriteString("# Chinese ipv4(*.*.*.0) got from http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest.")
	rd := bufio.NewReader(resp.Body)
	for {
		line ,e := rd.ReadString('\n')
		if e != nil{
			break
		}
		if strings.Contains(line ,"ipv4") && strings.Contains(line ,"CN") {
			f.WriteString(getIP(line))
			f.WriteString("\n")
		}
	}
}

// apnic|CN|ipv4|39.0.8.0|2048|20110412|allocated
func getIP(line string) string{
	nopre := line[14:]
	i := strings.Index(nopre ,"|")
	return nopre[:i-2]
}

// TODO not gentle
func loadCNIPList() error{
	f ,e := os.OpenFile(ipFile ,os.O_RDONLY ,0755)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	if line ,e := rd.ReadString('\n');e != nil{
		panic(e)
	}else{
		rmip := line[:len(line) - 2]
		cniplist[rmip] = struct{}{}
	}
	return nil
}


type domain struct {
	Bypass	map[string]uint8	`toml:"bypass"`
	Ads		map[string]uint8	`toml:"ads"`
}

func loadDomainList() domain{
	do := &domain{}
	if _ ,e := toml.DecodeFile(domainFile ,do);e != nil{
		panic(e)
	}
	return *do
}

const (
	BypassProxy = iota
	BypassTrue
	BypassDiscard
)
// Bypass 0:proxy ,1:bypass ,2:discard
func Bypass(dm string) int{
	for k := range domainlist.Bypass{
		if strings.HasSuffix(dm ,k) {
			return BypassTrue
		}
	}

	for k := range domainlist.Ads{
		if strings.HasSuffix(dm ,k){
			return BypassDiscard
		}
	}
	return BypassProxy
}

