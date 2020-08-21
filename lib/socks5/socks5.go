package socks5

import (
	"errors"
	"fmt"
	"io"
)

const (
	SocksVersion = 0x05

	SocksNoAuthenticationRequired = 0x00
	SocksUSERNAMEPASSWORD = 0x02

	SocksCmdConnect = 0x01

	SocksAddrTypeIPv4 = 0x01
	SocksAddrTypeDomain = 0x03
	SocksAddrTypeIPv6 = 0x04
)
// only support the method 0x00 now
var socks5HandshakeResp = []byte{5 ,0}

func Handshake(stream io.ReadWriter) error{
	buf := make([]byte ,255)
	if _ ,e := io.ReadFull(stream ,buf[:2]);e != nil{
		return e
	}
	if buf[0] != SocksVersion {
		return fmt.Errorf("unsupported version ,expect %d ,got %d" ,SocksVersion ,buf[0])
	}
	if buf[1] == 0{
		return fmt.Errorf("at least one supported authentication method")
	}
	if _, e := io.ReadFull(stream ,buf[:buf[1]]);e != nil{
		return e
	}
	// TODO full check ,error if 0x00 is not in first bit
	if buf[0] != SocksNoAuthenticationRequired{
		return errors.New("client must support the method [0x00 NO AUTHENTICATION REQUIRED]")
	}
	if _ ,e := stream.Write(socks5HandshakeResp);e != nil{
		return e
	}
	return nil
}

func ParseCMD(stream io.ReadWriter) (Addr ,error) {
	addr := Addr{}
	buf := make([]byte ,4)
	if _ ,e := io.ReadFull(stream ,buf[:4]);e != nil{
		return addr ,e
	}
	// skip version check
	if buf[1] != SocksCmdConnect{
		return addr ,errors.New("unsupported client cmd")
	}
	addr ,e := parseAddr(buf[3] ,stream)
	if e != nil{
		return addr ,e
	}

	// do something
	return addr ,nil
}

type Addr struct {
	IsDomain 	bool
	Port	 	[]byte
	IPOrDomain 	[]byte
}

func (addr Addr) String() string{
	port := int(addr.Port[0]) * 256 + int(addr.Port[1])
	if addr.IsDomain{
		return fmt.Sprintf("%s:%d" ,string(addr.IPOrDomain) ,port)
	}
	return fmt.Sprintf("%+v:%d" ,addr.IPOrDomain ,port)
}

func parseAddr(addrType byte ,stream io.ReadWriter) (Addr ,error){
	addr := Addr{}
	buf := make([]byte ,255)
	switch addrType {
	case SocksAddrTypeIPv4:
		if _ ,e := io.ReadFull(stream ,buf[:6]);e != nil{
			return addr ,e
		}
		addr.IPOrDomain = buf[:4]
		addr.Port = buf[4:6]
	case SocksAddrTypeIPv6:
		if _ ,e := io.ReadFull(stream ,buf[:18]);e != nil{
			return addr ,e
		}
		addr.IPOrDomain = buf[:16]
		addr.Port = buf[16:18]
	case SocksAddrTypeDomain:
		if _ ,e := io.ReadFull(stream ,buf[:1]);e != nil{
			return addr ,e
		}
		length := buf[0]
		if _ ,e := io.ReadFull(stream ,buf[:length + 2]);e != nil{
			return addr ,e
		}
		addr.IsDomain = true
		addr.IPOrDomain = buf[:length]
		addr.Port = buf[length:length+2]
	default:
		return addr ,errors.New("unknown address type")
	}
	return addr ,nil
}
