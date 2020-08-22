package transport

import (
	_ "golang.org/x/time/rate"
	"io"
)

const defaultBuflen = 1 << 14

func CopyAll(dst io.Writer ,src io.Reader ,buflen int) (int ,error){
	if buflen <= 0 {
		buflen = defaultBuflen
	}

	return 0 ,nil
}


