package main

import (
	"strings"
	"sync"

	"github.com/nynicg/cake/lib/log"
)

const ApnicUrl = "http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"

const (
	ipFile     = "apnic-latest.txt"
)

var domainCac domainCache

func init() {
	domainCac = domainCache{
		cache: map[string]int{},
	}
}

const (
	BypassProxy = iota
	BypassTrue
	BypassDiscard
)

// Bypass 0:proxy ,1:bypass ,2:discard
func Bypass(dm string) int {
	i, ok := GetDomainCache(dm)
	if ok {
		return i
	}

	for _, v := range cliConfig.RuleCfg.Bypass {
		if strings.HasSuffix(dm, v) {
			PutDomainCache(dm, BypassTrue)
			return BypassTrue
		}
	}

	for _, v := range cliConfig.RuleCfg.Discard {
		if strings.HasSuffix(dm, v) {
			PutDomainCache(dm, BypassDiscard)
			return BypassDiscard
		}
	}

	PutDomainCache(dm, BypassProxy)
	return BypassProxy
}

type domainCache struct {
	cache map[string]int
	mux   sync.Mutex
}

func PutDomainCache(domain string, rule int) {
	domainCac.mux.Lock()
	domainCac.cache[domain] = rule
	domainCac.mux.Unlock()
}

func GetDomainCache(domain string) (int, bool) {
	domainCac.mux.Lock()
	i, r := domainCac.cache[domain]
	domainCac.mux.Unlock()
	log.Debug("get domain bypass rule from cache ", domain, " -> ", i, " in cache ", r)
	return i, r
}
