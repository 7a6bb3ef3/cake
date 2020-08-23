package encrypt

import "testing"

func TestNewAES128CFB(t *testing.T) {
	cfb ,e := NewAES128CFB("1145141145141145")
	if e != nil{
		t.Fatal(e)
	}
	out ,e := cfb.Encrypt([]byte("123123123"))
	if e != nil{
		t.Fatal(e)
	}

	raw ,e := cfb.Decrypt(out)
	if e != nil{
		t.Fatal(e)
	}

	t.Log(string(raw))
}
