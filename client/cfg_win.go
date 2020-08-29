// +build windows

package main

import (
	"github.com/nynicg/cake/lib/log"
	"golang.org/x/sys/windows/registry"
)

// https://stackoverflow.com/questions/4283027/whats-the-format-of-the-defaultconnectionsettings-value-in-the-windows-registry
//keep this value.
//    1.  "00" placeholder
//    2.  "00" placeholder
//    3.  "00" placeholder
//    4.  "xx" increments if changed
//    5.  "xx" increments if 4. is "FF"
//    6.  "00" placeholder
//    7.  "00" placeholder
//    8.  "01"=proxy deaktivated; other value=proxy enabled
//    9.  "00" placeholder
//    10. "00" placeholder
//    11. "00" placeholder
//    12. "xx" length of "proxyserver:port"
//    13. "00" placeholder
//    14. "00" placeholder
//    15. "00" placeholder
//    "proxyserver:port"
//    if 'Bypass proxy for local addresses':::
//    other stuff with unknown length
//    "<local>"
//    36 times "00"
//    if no 'Bypass proxy for local addresses':::
//    40 times "00"
// TODO check the format ,actually it is different from the my local version(win10 18363) below
var defaultRegVal = []byte{
	70, 0, 0, 0, 50, 1, 0, 0, 3, 0, 0, 0, 14, 0, 0, 0, 49, 50, 55, 46, 48, 46, 48, 46, 49, 58, 49, 57, 49, 57, 182,
	0, 0, 0, 108, 111, 99, 97, 108, 104, 111, 115, 116, 59, 49, 50, 55, 46, 42, 59, 49, 48, 46, 42, 59, 49, 55, 50,
	46, 49, 54, 46, 42, 59, 49, 55, 50, 46, 49, 55, 46, 42, 59, 49, 55, 50, 46, 49, 56, 46, 42, 59, 49, 55, 50, 46,
	49, 57, 46, 42, 59, 49, 55, 50, 46, 50, 48, 46, 42, 59, 49, 55, 50, 46, 50, 49, 46, 42, 59, 49, 55, 50, 46, 50,
	50, 46, 42, 59, 49, 55, 50, 46, 50, 51, 46, 42, 59, 49, 55, 50, 46, 50, 52, 46, 42, 59, 49, 55, 50, 46, 50, 53,
	46, 42, 59, 49, 55, 50, 46, 50, 54, 46, 42, 59, 49, 55, 50, 46, 50, 55, 46, 42, 59, 49, 55, 50, 46, 50, 56, 46,
	42, 59, 49, 55, 50, 46, 50, 57, 46, 42, 59, 49, 55, 50, 46, 51, 48, 46, 42, 59, 49, 55, 50, 46, 51, 49, 46, 42,
	59, 49, 57, 50, 46, 49, 54, 56, 46, 42, 59, 60, 108, 111, 99, 97, 108, 62, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func getDefaultConnSetting() ([]byte, error) {
	key, _, e := registry.CreateKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Connections", registry.ALL_ACCESS)
	if e != nil {
		return nil, e
	}
	defer key.Close()
	raw, _, e := key.GetBinaryValue("DefaultConnectionSettings")
	if e != nil {
		return raw, e
	}
	if raw[4] == 255 {
		raw[4] = 0
	} else {
		raw[4] = raw[4] + 1
	}
	defaultRegVal[4] = raw[4]
	return defaultRegVal[:], nil
}

// configure modify the windows registry to enable the system proxy
func configure() {
	b, e := getDefaultConnSetting()
	if e != nil {
		log.Error(e)
		return
	}
	b[8] = 3
	key, _, e := registry.CreateKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Connections", registry.ALL_ACCESS)
	if e != nil {
		log.Error(e)
		return
	}
	defer key.Close()
	key.SetBinaryValue("DefaultConnectionSettings", b)
}

// unconfigure like configure ,modify the windows registry to disable the proxy
func unconfigure() {
	b, e := getDefaultConnSetting()
	if e != nil {
		log.Error(e)
		return
	}
	b[8] = 1
	key, _, e := registry.CreateKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Connections", registry.ALL_ACCESS)
	if e != nil {
		log.Error(e)
		return
	}
	defer key.Close()
	key.SetBinaryValue("DefaultConnectionSettings", b)
}
