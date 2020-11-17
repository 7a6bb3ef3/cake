package socks5

import (
	"errors"
	"fmt"
	"io"
)

// RFC1928
const (
	SocksVersion = 0x05

	SocksNoAuthenticationRequired = 0x00
	SocksUSERNAMEPASSWORD         = 0x02

	SocksCmdConnect = 0x01

	SocksAddrTypeIPv4   = 0x01
	SocksAddrTypeDomain = 0x03
	SocksAddrTypeIPv6   = 0x04

	SocksRespOK                = 0x00
	SocksRespServErr           = 0x01
	SocksRespHostUnreachable   = 0x04
	SocksRespUnsupportCmd      = 0x07
	SocksRespUnsupportAddrType = 0x08
)

// only support the method 0x00 now
var socks5HandshakeResp = []byte{5, 0}

func Handshake(stream io.ReadWriter) error {
	buf := make([]byte, 255)
	if _, e := io.ReadFull(stream, buf[:2]); e != nil {
		return e
	}
	if buf[0] != SocksVersion {
		return fmt.Errorf("unsupported version ,expect %d ,got %d", SocksVersion, buf[0])
	}
	if buf[1] == 0 {
		return fmt.Errorf("at least one supported authentication method")
	}
	if _, e := io.ReadFull(stream, buf[:buf[1]]); e != nil {
		return e
	}
	// TODO full check ,error if 0x00 is not in first bit
	if buf[0] != SocksNoAuthenticationRequired {
		return errors.New("client must support the method [0x00 NO AUTHENTICATION REQUIRED]")
	}
	if _, e := stream.Write(socks5HandshakeResp); e != nil {
		return e
	}
	return nil
}

func ParseCMD(stream io.ReadWriter) (Addr, error) {
	addr := Addr{}
	buf := make([]byte, 4)
	if _, e := io.ReadFull(stream, buf[:4]); e != nil {
		return addr, e
	}
	// skip version check
	if buf[1] != SocksCmdConnect {
		ProxyFailed(SocksRespUnsupportCmd, stream)
		return addr, errors.New("unsupported client cmd")
	}
	addr, e := parseAddr(buf[3], stream)
	if e != nil {
		return addr, e
	}

	// do something
	return addr, nil
}

func parseAddr(addrType byte, stream io.ReadWriter) (Addr, error) {
	addr := Addr{}
	buf := make([]byte, 255)
	switch addrType {
	case SocksAddrTypeIPv4:
		if _, e := io.ReadFull(stream, buf[:6]); e != nil {
			return addr, e
		}
		addr.Hostx = buf[:4]
		addr.Port = calcuPort(buf[4], buf[5])
	//case SocksAddrTypeIPv6:
	//if _ ,e := io.ReadFull(stream ,buf[:18]);e != nil{
	//	return addr ,e
	//}
	//addr.Domain = buf[:16]
	//addr.Port = calcuPort(buf[16] ,buf[16])
	case SocksAddrTypeDomain:
		if _, e := io.ReadFull(stream, buf[:1]); e != nil {
			return addr, e
		}
		length := buf[0]
		if _, e := io.ReadFull(stream, buf[:length+2]); e != nil {
			return addr, e
		}
		addr.IsDomain = true
		addr.Hostx = buf[:length]
		addr.Port = calcuPort(buf[length], buf[length+1])
	default:
		ProxyFailed(SocksRespUnsupportAddrType, stream)
		return addr, errors.New("unsupport address type")
	}
	return addr, nil
}

func calcuPort(a, b byte) int {
	return int(a)*256 + int(b)
}

func ProxyFailed(respCode byte, stream io.Writer) error {
	_, e := stream.Write([]byte{SocksVersion, respCode, 0, SocksAddrTypeIPv4, 0, 0, 0, 0, 0, 0})
	return e
}

func ProxyOKWithIpv4(stream io.Writer, ipv4 Addr) error {
	_, e := stream.Write([]byte{SocksVersion, SocksRespOK, 0, SocksAddrTypeIPv4, ipv4.Hostx[0],
		ipv4.Hostx[1], ipv4.Hostx[2], ipv4.Hostx[3], ipv4.BytePort()[0], ipv4.BytePort()[1],
	})
	return e
}

func ProxyOK(stream io.Writer) error {
	_, e := stream.Write([]byte{SocksVersion, SocksRespOK, 0, SocksAddrTypeIPv4, 0, 0, 0, 0, 0, 0})
	return e
}
