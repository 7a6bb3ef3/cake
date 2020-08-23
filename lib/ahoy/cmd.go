package ahoy

import (
	"errors"
)

type Command byte

const (
	CommandConnect Command = iota + 1
)


// customer remote cmd proto
//  +-----+-----+-----+-----+-----+
//  |ENC  |CMD  |ACKey|LEN  |ADDR |
//  +-----+-----+-----+-----+-----+
//  |1    |1    |16   |1    |LEN  |
// if success ,server response(random 6 bit):
//  +-----+-----+-----+-----+-----+-----+
//  |1    |1    |4    |5    |1    |4    |
type RemoteConnectRequest struct {
	Encryption	byte
	Command 	Command
	AccessKey	[]byte
	AddrLength	byte
	Addr		[]byte
}

func (r RemoteConnectRequest) Bytes() ([]byte ,error){
	if len(r.AccessKey) != 16 || len(r.Addr) != int(r.AddrLength) || int(r.AddrLength) == 0{
		return nil ,errors.New("invalid request")
	}
	buf := make([]byte ,1 + 1 + 16 + 1 + int(r.AddrLength))
	buf[0] = r.Encryption
	buf[1] = byte(r.Command)
	for i:=2;i<18;i++{
		buf[i] = r.AccessKey[i-2]
	}
	buf[18] = r.AddrLength
	for i:=19;i<19+int(r.AddrLength);i++{
		buf[i] = r.Addr[i-19]
	}
	return buf ,nil
}