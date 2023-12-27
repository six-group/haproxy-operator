package hash

import (
	"crypto/md5" //#nosec
	"encoding/hex"
)

func GetMD5Hash(values string) string {
	hash := md5.Sum([]byte(values)) //#nosec
	return hex.EncodeToString(hash[:])
}
