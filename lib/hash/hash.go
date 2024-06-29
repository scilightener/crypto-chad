package hash

import (
	"crypto/sha256"
)

func Hash(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}
