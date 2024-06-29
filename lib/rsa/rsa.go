package rsa

import (
	"crypto-chad-lib/mathf"
	"math/big"
)

const ChunkSizeBytes = 128

func Encrypt(message []byte, pow *big.Int, mod *big.Int) [][]byte {
	chunks := make([][]byte, 0, len(message)/ChunkSizeBytes+1)
	var i int
	for i = 0; i+ChunkSizeBytes < len(message); i += ChunkSizeBytes {
		chunks = append(chunks, transformBytes(message[i:i+ChunkSizeBytes], pow, mod))
	}
	if i < len(message) {
		chunks = append(chunks, transformBytes(message[i:], pow, mod))
	}

	return chunks
}

func Decrypt(chunks [][]byte, pow *big.Int, mod *big.Int) []byte {
	res := make([]byte, 0)
	for i := range chunks {
		res = append(res, transformBytes(chunks[i], pow, mod)...)
	}

	return res
}

func transformBytes(b []byte, pow, mod *big.Int) []byte {
	num := new(big.Int).SetBytes(b)
	if num.Cmp(big.NewInt(0)) <= 0 {
		panic("rsa.transformBytes: pow should be > 0")
	}
	if pow.Cmp(big.NewInt(0)) <= 0 {
		panic("rsa.transformBytes: pow should be > 0")
	}
	if mod.Cmp(big.NewInt(0)) <= 0 {
		panic("rsa.transformBytes: pow should be > 0")
	}
	exp := mathf.PowerMod(num, pow, mod)
	return exp.Bytes()
}
