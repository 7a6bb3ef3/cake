package encrypt

var defaultPlain = &Plain{}

func GetTypePlain() Encryptor{
	return defaultPlain
}

type Plain struct {}

func (p Plain) Encrypt(in []byte) (out []byte, err error) {
	return in ,nil
}

func (p Plain) Decrypt(in []byte) (out []byte, err error) {
	return in ,nil
}

