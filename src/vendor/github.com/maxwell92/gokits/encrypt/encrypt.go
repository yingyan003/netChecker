package encrypt

import (
	"crypto/sha256"
	"encoding/base64"
)

const (
	ENCRYPT_SALT1 = "999983"
	ENCRYPT_SALT2 = "22c362"
)

type Encryption struct {
	password string
}

func NewEncryption(password string) *Encryption {
	encrypt := ENCRYPT_SALT1 + password + ENCRYPT_SALT2
	sha := sha256.New()
	sha.Write([]byte(encrypt))
	pass := base64.StdEncoding.EncodeToString(sha.Sum(nil))
	return &Encryption{password: pass}
}

func (e *Encryption) String() string {
	return e.password
}
