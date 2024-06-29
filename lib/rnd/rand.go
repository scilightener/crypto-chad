package rnd

import (
	"math/rand/v2"
	"strings"
)

const alph = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func String(length int) string {
	var sb strings.Builder
	for range length {
		l := rand.IntN(len(alph))
		sb.WriteByte(alph[l])
	}

	return sb.String()
}
