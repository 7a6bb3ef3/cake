package ahoy

import (
	"errors"
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/mathx"
	"io"
	"net"
)

// use a customer protocol ,for experiment
// return encryption type ,proxy address and an error if there is
func Handshake(ackey string ,fromsocks net.Conn) (int ,string ,error){
	buf := make([]byte ,255)
	if _ ,e := io.ReadFull(fromsocks ,buf[:19]);e != nil{
		return 0 ,"" ,e
	}
	log.Debug("handshake pack " ,buf[:19])
	addrLen := buf[18]
	enctype := buf[0]
	if buf[1] != byte(CommandConnect) {
		return 0 ,"" ,errors.New("unsupport command")
	}else if string(buf[2:18]) != ackey {
		return 0 ,"" ,errors.New("access refused")
	}else if addrLen == 0{
		return 0 ,"" ,errors.New("empty proxy addr")
	}
	// read addr
	if _ ,e := io.ReadFull(fromsocks ,buf[:addrLen]);e != nil{
		return 0 ,"" ,e
	}
	return int(enctype) ,string(buf[:addrLen]) ,nil
}

func OnReady(w io.Writer) error{
	_ ,e := w.Write([]byte{mathx.Byten(255) ,mathx.Byten(255) ,mathx.Byten(255) ,mathx.Byten(255) ,mathx.Byten(255) ,mathx.Byten(255)})
	return e
}
