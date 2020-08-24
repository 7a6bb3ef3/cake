package ahoy

import (
	"io"

	"github.com/nynicg/cake/lib/encrypt"
	"github.com/nynicg/cake/lib/log"
)

const blockSize = 1024


// CopyWithCryptFunc both encryption and decryption func is welcome here ,
// anyway they are used to process data read from src,then write to dst
func CopyWithCryptFunc(dst io.Writer ,src io.Reader ,encrypt encrypt.EncryptFunc ,buf []byte) (int ,error){
	var (
		written = 0
		err error
	)
	// https://github.com/golang/go/blob/master/src/io/io.go
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			encrypted ,e := encrypt(buf[0:nr])
			if e != nil{
				log.Error("encrypted/decrypted in copy." ,e ," length of bytes " ,nr)
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

//// CopyWithDecryptFunc
//func CopyWithDecryptFunc(dst io.Writer ,src io.Reader ,decrypt encrypt.EncryptFunc ,buf []byte) (int ,error){
//	var (
//		written = 0
//		err error
//		left = []byte{}
//		toDecrypt []byte
//	)
//	// https://github.com/golang/go/blob/master/src/io/io.go
//	for {
//		nr, er := src.Read(buf)
//		if nr > 0 {
//			index ,ok := findEnd(buf[:nr])
//			if ok {
//				toDecrypt = append(left ,buf[:index]...)
//				if index + 8 < nr {
//					left = append(left ,buf[index+8:]...)
//				}
//			}else{
//				left = append(left ,buf[:nr]...)
//				continue
//			}
//			encrypted ,e := decrypt(toDecrypt)
//			if e != nil{
//				log.Error("encrypted/decrypted in copy." ,e ," length of bytes " ,nr)
//				return written ,e
//			}
//			out := append(encrypted ,getEnd()...)
//			nw, ew := dst.Write(out)
//			if nw > 0 {
//				written += nw
//			}
//			if ew != nil {
//				err = ew
//				break
//			}
//		}
//		if er != nil {
//			if er != io.EOF {
//				err = er
//			}
//			break
//		}
//	}
//	return written ,err
//}
//
//
//func getEnd() []byte{
//	utc := time.Now().UTC()
//	msg := fmt.Sprintf("19%d89%dyjSp%sSrmy%d" ,utc.Day() ,utc.Year() ,utc.Month().String(),utc.Month())
//	sum := sha256.Sum256([]byte(msg))
//	return []byte{sum[utc.Day()] ,sum[6] ,sum[utc.Month()*2] ,sum[8] ,sum[utc.Month()] ,sum[3] ,sum[17] ,sum[11]}
//}
//
//func findEnd(in []byte) (int ,bool) {
//	if i := bytes.Index(in ,getEnd());i == -1{
//		return 0 ,false
//	}else{
//		return i ,true
//	}
//	//return i ,i > -1
//}
