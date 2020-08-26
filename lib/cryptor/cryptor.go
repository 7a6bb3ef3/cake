package cryptor

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"sync"
)

const (
	CryptTypeAES128GCM = iota + 1
	// xchacha20poly1305
	CryptTypeCHACHA
	CryptTypePlain
)

const (
	NameAES128GCM = "aes128gcm"
	NameCHACHA = "chacha"
	NamePlain = "plain"
)

// both encryption or decryption func is ok
type CryptFunc func(in []byte) (out []byte ,e error)

// StreamEncryptor used to cryptor communication in net.conn.
type Cryptor interface {
	Encrypt(in []byte) (out []byte ,err error)
	Decrypt(in []byte) (out []byte ,err error)
}

func GetStreamEncryptorIndexByName(name string) (int ,error){
	switch strings.ToLower(name) {
	case NameAES128GCM:
		return CryptTypeAES128GCM ,nil
	case NameCHACHA:
		return CryptTypeCHACHA ,nil
	case NamePlain:
		return CryptTypePlain ,nil
	default:
		return 0 ,fmt.Errorf("no such encryptor %s" ,name)
	}
}


type CryptorMap struct {
	m 	map[int]Cryptor
	mutex sync.Mutex
}

func NewEncryptorMap() *CryptorMap{
	return &CryptorMap{
		m: make(map[int]Cryptor),
	}
}

func (e *CryptorMap) Get(index int) (Cryptor ,error){
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if en ,ok := e.m[index];!ok {
		return nil ,errors.New("unregistered encryptor")
	}else{
		return en ,nil
	}
}

func (e *CryptorMap) Register(index int ,en Cryptor) error{
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _ ,ok := e.m[index];ok {
		return fmt.Errorf("already registered index number %d" ,index)
	}
	e.m[index] = en
	return nil
}

// retrun result[:n]
func sha256N(s string ,n int) []byte{
	re := sha256.Sum256([]byte(s))
	rex := re[:]
	// out bounds check
	return rex[:n]
}