package cryptor

var defaultPlain = &Plain{}

func GetTypePlain() Cryptor{
	return defaultPlain
}

type Plain struct {}

func (p Plain) Encrypt(in []byte) (out []byte, err error) {
	return in ,nil
}

func (p Plain) Decrypt(in []byte) (out []byte, err error) {
	return in ,nil
}

