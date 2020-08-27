package ahoy

import "testing"

func TestRemoteProxyRequest_Bytes(t *testing.T) {
	req := RemoteConnectRequest{
		Encryption: 0,
		Command:    0,
		RandomKey: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		AddrLength: 14,
		Addr:       []byte("google.com:443"),
	}
	b, e := req.Bytes()
	if e != nil {
		t.Fatal(e)
	}
	t.Log(b)
	t.Log(req.Addr)
}
