package socks5

import "testing"

func TestAddr_BytePort(t *testing.T) {
	addr := Addr{
		Port: 443,
	}
	t.Log(addr.BytePort())
	addr.Port = 80
	t.Log(addr.BytePort())
	addr.Port = 8080
	t.Log(addr.BytePort())
}
