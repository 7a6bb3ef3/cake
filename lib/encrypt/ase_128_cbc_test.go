package encrypt

import (
	"bytes"
	"log"
	"testing"
)

func TestAES128CBC_EncryptStream(t *testing.T) {
	aes ,e := NewAES128CBC("1145141145141145" ,"1145141145141145")
	if e != nil{
		t.Fatal(e)
	}

	in := bytes.NewBuffer([]byte{})
	in.WriteString("simple msg")
	dst1 := bytes.NewBuffer([]byte{})
	if e := aes.EncryptStream(dst1 ,in);e != nil{
		t.Fatal(e)
	}
	dst2 := bytes.NewBuffer([]byte{})
	t.Log(dst1.String())
	if e := aes.DecryptStream(dst2 ,dst1);e != nil{
		t.Fatal(e)
	}
	t.Log(dst2.String())
}

func TestAES128CBC_Encrypt(t *testing.T) {
	aes ,e := NewAES128CBC("1145141145141145" ,"1145141145141145")
	if e != nil{
		t.Fatal(e)
	}
	out ,e := aes.Encrypt([]byte("simple msg"))
	if e != nil{
		t.Fatal(e)
	}
	t.Log(out)



	aes2 ,e := NewAES128CBC("1145141145141145" ,"1145141145141145")
	if e != nil{
		t.Fatal(e)
	}
	raw2 ,e := aes2.Decrypt(out)
	if e != nil{
		t.Fatal(e)
	}
	log.Println(string(raw2))
}


// BenchmarkAES128CBC_Encrypt-8
//  2349862 504 ns/op 336 B/op 12 allocs/op
func BenchmarkAES128CBC_Encrypt(b *testing.B) {
	aes ,e := NewAES128CBC("1145141145141145" ,"1145141145141145")
	if e != nil{
		b.Fatal(e)
	}
	in := []byte("simple msg")
	b.ReportAllocs()
	for i:=0;i<b.N*2;i++ {
		_ ,e = aes.Encrypt(in)
		if e != nil{
			b.Fatal(e)
		}
	}
}


// BenchmarkAES128CBC_Decrypt
//  3602426 333 ns/op 256 B/op 8 allocs/op
func BenchmarkAES128CBC_Decrypt(b *testing.B) {
	aes ,e := NewAES128CBC("1145141145141145" ,"1145141145141145")
	if e != nil{
		b.Fatal(e)
	}
	b.ReportAllocs()
	in := []byte{93 ,249 ,39 ,135 ,47 ,2 ,76 ,105 ,157 ,67 ,108 ,11 ,130 ,107 ,111 ,222}
	for i:=0;i<b.N*2;i++ {
		_ ,e = aes.Decrypt(in)
		if e != nil{
			b.Fatal(e)
		}
	}
}
