package cryptor

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	NameCHACHA    = "chacha"
	NamePlain     = "plain"
)

// both encryption or decryption func is ok
type CryptFunc func(in []byte) (out []byte, e error)

// NewCryptorX return random key
func NewCryptorX(i int) (Cryptor, string, error) {
	k := random32key()
	cr, e := NewCryptor(i, k)
	if e != nil {
		return nil, k, e
	}
	return cr, k, nil
}

func random32key() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func NewCryptor(i int, key string) (Cryptor, error) {
	switch i {
	case CryptTypeAES128GCM:
		return NewAES128GCM(key)
	case CryptTypeCHACHA:
		return NewChacha20Poly1305(key)
	case CryptTypePlain:
		return &Plain{}, nil
	}
	return nil, errors.New("unsupported cryption")
}

// StreamEncryptor used to cryptor communication in net.conn.
type Cryptor interface {
	Encrypt(in []byte) (out []byte, err error)
	Decrypt(in []byte) (out []byte, err error)
}

func GetIndexByName(name string) (int, error) {
	switch strings.ToLower(name) {
	case NameAES128GCM:
		return CryptTypeAES128GCM, nil
	case NameCHACHA:
		return CryptTypeCHACHA, nil
	case NamePlain:
		return CryptTypePlain, nil
	default:
		return 0, fmt.Errorf("no such encryptor %s", name)
	}
}

type CryptorMap struct {
	m     map[int]Cryptor
	mutex sync.Mutex
}

func NewEncryptorMap() *CryptorMap {
	return &CryptorMap{
		m: make(map[int]Cryptor),
	}
}

func (e *CryptorMap) Get(index int) (Cryptor, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if en, ok := e.m[index]; !ok {
		return nil, errors.New("unregistered encryptor")
	} else {
		return en, nil
	}
}

func (e *CryptorMap) Register(index int, en Cryptor) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _, ok := e.m[index]; ok {
		return fmt.Errorf("already registered index number %d", index)
	}
	e.m[index] = en
	return nil
}

func RegistryAllCrypto(key string) *CryptorMap {
	enmap := NewEncryptorMap()

	enmap.Register(CryptTypePlain, &Plain{})

	gcm, e := NewAES128GCM(key)
	if e != nil {
		panic(e)
	}
	enmap.Register(CryptTypeAES128GCM, gcm)

	cc, e := NewChacha20Poly1305(key)
	if e != nil {
		panic(e)
	}
	enmap.Register(CryptTypeCHACHA, cc)
	return enmap
}

// retrun result[:n]
func sha256N(s string, n int) []byte {
	re := sha256.Sum256([]byte(s))
	rex := re[:]
	// out bounds check
	return rex[:n]
}
