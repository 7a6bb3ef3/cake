package ahoy

import (
	"errors"
	"github.com/nynicg/cake/lib/mathx"
	"io"
	"net"
)

// use a customer protocol ,for experiment
func Handshake(ackey string ,fromsocks net.Conn) (string ,error){
	buf := make([]byte ,255)
	if _ ,e := io.ReadFull(fromsocks ,buf[:19]);e != nil{
		return "" ,e
	}
	addrLen := buf[18]
	if buf[1] != byte(CommandConnect) {
		return "" ,errors.New("unsupport command")
	}else if string(buf[2:18]) != ackey {
		return "" ,errors.New("access refused")
	}else if addrLen == 0{
		return "" ,errors.New("empty proxy addr")
	}
	// read addr
	if _ ,e := io.ReadFull(fromsocks ,buf[:addrLen]);e != nil{
		return "" ,e
	}
	return string(buf[:addrLen]) ,nil
}

func OnReady(w io.Writer) error{
	_ ,e := w.Write([]byte{mathx.Byten(255) ,mathx.Byten(255) ,mathx.Byten(255) ,mathx.Byten(255) ,mathx.Byten(255) ,mathx.Byten(255)})
	return e
}
