package tool

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encrypt(start string) string {
	h := md5.New()
	h.Write([]byte(start))
	return hex.EncodeToString(h.Sum(nil))
}
