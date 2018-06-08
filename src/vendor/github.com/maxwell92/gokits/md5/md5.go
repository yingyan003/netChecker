package common

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Sum(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
