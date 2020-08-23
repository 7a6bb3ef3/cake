package encrypt

import (
	"io"
)

// StreamEncryptor used to encrypt communication in net.conn.
type StreamEncryptor interface {
	// EncryptStream encrypt the plain data from in and write to out ,
	// out will get encrypted data
	EncryptStream(dst io.Writer ,src io.Reader) error
	// EncryptStream decrypt the encryted data from in and write to out ,
	// out will get the plain data
	DecryptStream(dst io.Writer ,src io.Reader) error

	Encrypt(in []byte) (out []byte ,err error)
	Decrypt(in []byte) (out []byte ,err error)
}


// PiplineEncryptor a pipline style api is better?
//  var dst ,src net.Conn
//  ...
//  ip.Copy(dst ,pip.EncryptSteam(src))
// Deprecated it's inconvenient to handle error
type PiplineEncryptor interface {
	EncryptStream(in io.Reader) io.Reader
	DecryptStream(in io.Reader) io.Reader
}
