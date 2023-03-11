package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

func Compute(pass string) string {
	passSha := sha1.Sum([]byte(pass))
	passwordHash := hex.EncodeToString(passSha[:])
	return passwordHash
}
