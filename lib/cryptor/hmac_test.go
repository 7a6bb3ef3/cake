package cryptor

import (
	"math/rand"
	"testing"
)

func TestNewHMAC(t *testing.T) {
	in := buildInput()
	hm := NewHMAC("M5Rm2nmNyn1cg")
	for _ ,v := range in {
		out,e := hm.SumN([]byte(v) ,rand.Intn(32))
		out2 := hm.Sum([]byte(v))
		if e != nil{
			t.Error(e)
		}else{
			t.Log(string(out))
			t.Log("length" ,len(out2))
		}
	}
}
