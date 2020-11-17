package socks5

import "fmt"

// domain ,ipv4 supoort
type Addr struct {
	IsDomain   bool
	Port       int
	Hostx 	   []byte
}

// Address return addr ,e.g. 127.0.0.1:1919 ,google.com:443
func (addr Addr) Address() string {
	return fmt.Sprintf("%s:%d", addr.Host(), addr.Port)
}

// Host return domain in string or ipv4 host
func (addr Addr) Host() string {
	if addr.IsDomain {
		return string(addr.Hostx)
	}
	return fmt.Sprintf("%d.%d.%d.%d", addr.Hostx[0], addr.Hostx[1], addr.Hostx[2], addr.Hostx[3])
}

func (addr Addr) BytePort() []byte {
	f := addr.Port / 256
	s := addr.Port - f*256
	return []byte{byte(f), byte(s)}
}
