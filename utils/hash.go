package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func Hash[T string | []byte](buf T) string {
	hash := sha512.Sum512_256([]byte(buf))
	return hex.EncodeToString(hash[:])
}
