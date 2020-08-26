package ahoy

import (
	"errors"
	"fmt"
	"github.com/nynicg/cake/lib/cryptor"
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/pool"
	"io"
)

type CopyEnv struct {
	ReaderWithLength bool
	WriterNeedLength bool
	CryptFunc        cryptor.CryptFunc
	BufPool          *pool.BufferPool
	Bypass           bool
}


func CopyConn(dst io.Writer ,src io.Reader ,cfg *CopyEnv) (int ,error){
	buf := cfg.BufPool.Get()
	defer cfg.BufPool.Put(buf)
	if cfg.Bypass {
		cfg.ReaderWithLength = false
		cfg.WriterNeedLength = false
	}
	var (
		written int
		err error
		srcpayload []byte
	)
	for {
		if !cfg.ReaderWithLength{
			nr, er := src.Read(buf)
			err = er
			srcpayload = buf[:nr]
		}else{
			d ,e := readWithLength(src)
			err = e
			srcpayload = d
		}
		if err != nil && !errors.Is(err ,io.EOF){
			break
		}
		// no more payload to handle
		if errors.Is(err ,io.EOF) && len(srcpayload) == 0 {
			err = nil
			break
		}
		// got EOF and some payload ,or no error here
		towrite ,e := cfg.CryptFunc(srcpayload)
		if e != nil{
			return written ,fmt.Errorf("CopyConn.tcp read:%w ,crpyto:%s" ,err ,e.Error())
		}
		w ,e := writeWithLength(dst ,towrite ,cfg.WriterNeedLength)
		if e != nil{
			return written ,fmt.Errorf("CopyConn.tcp write:%w" ,e)
		}
		written += w

		// payload has been written.if EOF exist ,do break.
		if errors.Is(err ,io.EOF) {
			err = nil
			break
		}
	}
	return written ,err
}



// big-endian
func writeWithLength(writer io.Writer ,bytes []byte ,needLength bool) (int ,error){
	l := len(bytes)
	s := byte(l % 256)
	f := byte((l - int(s)) / 256 )
	written := 0
	if needLength {
		if _ ,e := writer.Write([]byte{f ,s});e != nil{
			return 0 ,fmt.Errorf("writeWithLength:%w" ,e)
		}
		written = 2
		log.Debug("write length head {" ,f ,s ,"}")
	}

	if n ,e := writer.Write(bytes);e != nil{
		return 0 ,fmt.Errorf("writeWithLength:%w" ,e)
	}else{
		log.Debug("finish write " ,n)
		return written+n ,nil
	}
}

// readWithLength
func readWithLength(rd io.Reader) ([]byte , error) {
	var (
		length int
		out []byte
	)

	lenBit := make([]byte ,2)
	_ ,e := io.ReadFull(rd ,lenBit)
	if e != nil{
		return nil ,fmt.Errorf("readWithLength:%w" ,e)
	}
	length = int(lenBit[0]) * 256 + int(lenBit[1])
	log.Debug("read pack has length head " ,length ,"bits")
	out = make([]byte ,length)
	_ ,e = io.ReadFull(rd ,out)
	if e != nil{
		return nil ,fmt.Errorf("readWithLength:%w" ,e)
	}
	log.Debug("finish read " ,len(out))
	return out ,nil
}