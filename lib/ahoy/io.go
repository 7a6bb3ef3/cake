package ahoy

import (
	"fmt"
	"io"

	"github.com/nynicg/cake/lib/cryptor"
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/pool"
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
		eof bool
	)
	for {
		if !cfg.ReaderWithLength{
			nr, er := src.Read(buf)
			if er != nil && er != io.EOF {
				err = er
				break
			}else if er == io.EOF{
				eof = true
			}
			srcpayload = buf[:nr]
			log.Debug("read no head pack -> " ,nr ," bits")
		}else{
			d ,e := readWithLength(src)
			if e != nil{
				err = e
				break
			}
			srcpayload = d
		}
		towrite ,e := cfg.CryptFunc(srcpayload)
		if e != nil{
			err = e
			break
		}
		w ,e := writeWithLength(dst ,towrite ,cfg.WriterNeedLength)
		if e != nil{
			err = e
			break
		}
		written += w
		if eof {
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
