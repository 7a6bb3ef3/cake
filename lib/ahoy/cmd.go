package ahoy

import (
	"errors"
)

type Command byte

const (
	HMACOK = iota
	HMACInvalid
)

const (
	CommandConnect Command = iota + 1
)

// customer remote cmd proto
//  +-----+-----+-----+-----+-----+
//  |ENC  |CMD  |RDKey|LEN  |ADDR |
//  +-----+-----+-----+-----+-----+
//  |1    |1    |32   |1    |LEN  |
// if success ,server response(random 6 bit)
// TODO take some specific data as resp ,like 1=ok ,2=failed or something
type RemoteConnectRequest struct {
	Encryption byte
	Command    Command
	RandomKey  []byte
	AddrLength byte
	Addr       []byte
}

func (r RemoteConnectRequest) Bytes() ([]byte, error) {
	if len(r.RandomKey) != 32 || len(r.Addr) != int(r.AddrLength) || int(r.AddrLength) == 0 {
		return nil, errors.New("invalid request")
	}
	buf := make([]byte, 1+1+32+1+int(r.AddrLength))
	buf[0] = r.Encryption
	buf[1] = byte(r.Command)
	for i := 2; i < 34; i++ {
		buf[i] = r.RandomKey[i-2]
	}
	buf[34] = r.AddrLength
	for i := 35; i < 35+int(r.AddrLength); i++ {
		buf[i] = r.Addr[i-35]
	}
	return buf, nil
}
