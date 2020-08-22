package mathx

import (
	"math/rand"
	"time"
)

var rd *rand.Rand

func init() {
	src := rand.NewSource(time.Now().UnixNano() + int64(rand.Intn(1 << 20)))
	rd = rand.New(src)
}

func Intn(n int) int {
	return rd.Intn(n)
}

func Byten(n int) byte{
	return byte(Intn(n))
}
