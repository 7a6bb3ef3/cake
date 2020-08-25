package encrypt

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"sync"
)

const (
	EncryptTypeAES128GCM = iota + 1
	// xchacha20poly1305
	EncryptTypeCHACHA
	EncryptTypePlain
)

// EncryptFunc both encryption or decryption func is ok
type EncryptFunc func(in []byte) (out []byte ,e error)

// StreamEncryptor used to encrypt communication in net.conn.
type Encryptor interface {
	Encrypt(in []byte) (out []byte ,err error)
	Decrypt(in []byte) (out []byte ,err error)
}

func GetStreamEncryptorIndexByName(name string) (int ,error){
	switch strings.ToLower(name) {
	case "aes128gcm":
		return EncryptTypeAES128GCM ,nil
	case "chacha":
		return EncryptTypeCHACHA ,nil
	case "plain":
		return EncryptTypePlain ,nil
	default:
		return 0 ,fmt.Errorf("no such encryptor %s" ,name)
	}
}


type EncryptorMap struct {
	m 	map[int]Encryptor
	mutex sync.Mutex
}

func NewEncryptorMap() *EncryptorMap{
	return &EncryptorMap{
		m: make(map[int]Encryptor),
	}
}

func (e *EncryptorMap) Get(index int) (Encryptor ,error){
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if en ,ok := e.m[index];!ok {
		return nil ,errors.New("unregistered encryptor")
	}else{
		return en ,nil
	}
}

func (e *EncryptorMap) Register(index int ,en Encryptor) error{
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