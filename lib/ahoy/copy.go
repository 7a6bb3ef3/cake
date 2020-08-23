package ahoy

import (
	"io"

	"github.com/nynicg/cake/lib/encrypt"
	"github.com/nynicg/cake/lib/log"
	"github.com/nynicg/cake/lib/pool"
)



var bufPool *pool.BufferPool

func init(){
	bufPool = pool.NewBufPool(32 * 1024)
}

func CopyWithEncryptor(dst io.Writer ,src io.Reader ,encrypt encrypt.EncryptFunc) (int ,error){
	buf := bufPool.Get()
	defer func() {
		bufPool.Put(buf)
		log.Debug("finish write encrypted stream")
	}()
	var (
		written = 0
		err error
	)
	// Go\src\io\io.go
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			encrypted ,e := encrypt(buf[0:nr])
			if e != nil{
				return written ,e
			}
			nw, ew := dst.Write(encrypted)
			if nw > 0 {
				written += nw
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written ,err
}
