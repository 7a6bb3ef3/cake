package socks5

import "fmt"

// domain ,ipv4 supoort
type Addr struct {
	IsDomain 	bool
	Port	 	int
	IPOrDomain 	[]byte
}

// Address return addr ,e.g. 127.0.0.1:1919 ,google.com:443
func (addr Addr) Address() string{
	return fmt.Sprintf("%s:%d" ,addr.Host() ,addr.Port)
}


// Host return domain in string or ipv4 host
func (addr Addr) Host() string{
	if addr.IsDomain {
		return string(addr.IPOrDomain)
	}
	return fmt.Sprintf("%d.%d.%d.%d" ,addr.IPOrDomain[0],addr.IPOrDomain[1],addr.IPOrDomain[2],addr.IPOrDomain[3])
}

func (addr Addr) BytePort() []byte{
	f := addr.Port / 256
	s := addr.Port - f * 256
	return []byte{byte(f) ,byte(s)}
}
