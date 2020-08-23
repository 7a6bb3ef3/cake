package encrypt

import (
	"errors"
	"io"
	"strings"
)

const (
	EncryptTypeAES128CBC = iota + 1
	EncryptTypePlain
)

// EncryptFunc both encryption or decryption func is ok
type EncryptFunc func(in []byte) (out []byte ,e error)

// StreamEncryptor used to encrypt communication in net.conn.
type Encryptor interface {
	Encrypt(in []byte) (out []byte ,err error)
	Decrypt(in []byte) (out []byte ,err error)
}

func GetStreamEncryptor(index int) (Encryptor ,error){
	switch index {
	case EncryptTypeAES128CBC:
		return defaultAES128CBC ,nil
	case EncryptTypePlain:
		return defaultPlain ,nil
	default:
		return nil ,errors.New("unknown encrytor type")
	}
}

func GetStreamEncryptorIndexByName(name string) int{
	switch strings.ToLower(name) {
	case "aes128cbc":
		return EncryptTypeAES128CBC
	case "plain":
		return EncryptTypePlain
	default:
		// TODO rm panic
		panic("no such stream encryptor")
	}
}

func GetStreamEncryptorByName(name string) (Encryptor ,error){
	switch strings.ToLower(name) {
	case "aes128cbc":
		return defaultAES128CBC ,nil
	case "plain":
		return defaultPlain ,nil
	default:
		return nil ,errors.New("unknown encrytor type")
	}
}


func NewStreamEncryptorByName(name ,key ,vi string) (Encryptor ,error){
	switch strings.ToLower(name) {
	case "aes128cbc":
		return NewAES128CBC(key ,vi)
	case "plain":
		return defaultPlain ,nil
	default:
		return nil ,errors.New("unknown encrytor type")
	}
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
