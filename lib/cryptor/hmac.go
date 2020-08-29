package cryptor

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"math/rand"
	"strconv"
	"time"
)

var rd *rand.Rand

func init(){
	src := rand.NewSource(time.Now().UnixNano())
	rd = rand.New(src)
}

func HMAC(uid string) []byte {
	uidbyte := []byte(uid)
	hs := hmac.New(md5.New, uidbyte)
	nums := fbMultipleten()
	salt := salt(nums[rd.Intn(5)])
	hs.Write(append([]byte(salt), uidbyte...))
	return hs.Sum(nil)
}

func HMACAllTime(uid string) [][]byte{
	nums := fbMultipleten()
	out := make([][]byte ,len(nums))
	for i ,v := range nums{
		uidbyte := []byte(uid)
		hs := hmac.New(md5.New, uidbyte)
		salt := salt(v)
		hs.Write(append([]byte(salt), uidbyte...))
		out[i] = hs.Sum(nil)
	}
	return out
}

func VerifyHMAC(uid string ,b []byte) bool{
	ns := fbMultipleten()
	uy := []byte(uid)
	for _ ,v := range ns{
		hs := hmac.New(md5.New, uy)
		sa := salt(v)
		hs.Write(append([]byte(sa) ,uy...))
		if bytes.Equal(b ,hs.Sum(nil)) {
			return true
		}
	}
	return false
}

func salt(i int64) string{
	return "korone_godDoggo" + strconv.Itoa(int(i))
}


// the number nearest to the timestamp which is multiple of 8.
func fbMultipleten() []int64{
	nums := make([]int64 ,8)
	cur := time.Now().Unix()
	offset := cur % 8

	nums[0] = (cur - offset) + 2 * 8
	nums[1] = (cur - offset) + 8
	nums[2] = (cur - offset) - 8
	nums[3] = cur - offset
	nums[4] = (cur - offset) - 2 * 8

	nums[5] = (cur - offset) + 3 * 8
	nums[6] = (cur - offset) - 3 * 8
	return nums
}