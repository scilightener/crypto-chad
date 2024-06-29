package common

import (
	"crypto-chad-lib/rsa"
	"math/big"
)

var (
	ServerPubKey *rsa.PublicKey
)

func SetServerPubKey(e, n string) {
	mod, ok := new(big.Int).SetString(n, 10)
	if !ok {
		panic("failed to set server public key n param")
	}
	exp, ok := new(big.Int).SetString(e, 10)
	if !ok {
		panic("failed to set server public key e param")
	}
	ServerPubKey = &rsa.PublicKey{
		E: exp,
		N: mod,
	}
}

func ValidateSign(message string, sign [][]byte) bool {
	return eq([]byte(message), rsa.Decrypt(sign, ServerPubKey.E, ServerPubKey.N))
}

func eq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
