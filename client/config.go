package main

import (
	"github.com/BurntSushi/toml"
)

const fn = "./client_configure.toml"

var globCfg *Configure

func init() {
	globCfg = loadConfig()
}

type Configure struct {
	ProxyCfg ProxyCfg `toml:"proxy"`
	RuleCfg  RuleCfg  `toml:"rule"`
}

type ProxyCfg struct {
	ServerAddr     []string `toml:"serverAddr"`
	ServerPerfer   string   `toml:"serverPerfer"`
	LocalSocksAddr string   `toml:"localSocksAddr"`
	LocalHttpAddr  string   `toml:"localHttpAddr"`
	LogLevel       string   `toml:"logLevel"`
	EncryptType    string   `toml:"encryptType"`
	Key            string   `toml:"key"`
	Uid            string   `toml:"uid"`
	DisableGui     bool     `toml:"disableGui"`
	EnforceProxy   bool     `toml:"enforceBypass"`
}

type RuleCfg struct {
	Bypass  []string `toml:"bypass"`
	Discard []string `toml:"discard"`
}

func loadConfig() *Configure {
	cfg := &Configure{}
	if _, e := toml.DecodeFile(fn, cfg); e != nil {
		panic(e)
	}
	return cfg
}

func overrideByCmd(dst *Configure, src *ProxyCfg) {
	if src.EnforceProxy {
		dst.ProxyCfg.EnforceProxy = true
	}
	if src.DisableGui {
		dst.ProxyCfg.DisableGui = true
	}
	if src.ServerPerfer != "" {
		dst.ProxyCfg.ServerPerfer = src.ServerPerfer
	}
	if src.LocalHttpAddr != "" {
		dst.ProxyCfg.LocalHttpAddr = src.LocalHttpAddr
	}
	if src.LocalSocksAddr != "" {
		dst.ProxyCfg.LocalSocksAddr = src.LocalSocksAddr
	}
	if src.EncryptType != "" {
		dst.ProxyCfg.EncryptType = src.EncryptType
	}
	if src.Key != "" {
		dst.ProxyCfg.Key = src.Key
	}
	if src.Uid != "" {
		dst.ProxyCfg.Uid = src.Uid
	}
	if src.LogLevel != "" {
		dst.ProxyCfg.LogLevel = src.LogLevel
	}
}
