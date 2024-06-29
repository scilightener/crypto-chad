package ecdh

import (
	"crypto-chad-lib/ec"
	"math/big"
)

var DefaultECDH = createDefaultECDH()

func createDefaultECDH() *ECDH {
	curve := ec.DefaultCurve
	g := createDefaultG()
	n := createDefaultN()
	return NewECDH(*curve, g, n)
}

func createDefaultG() *ec.Point {
	x, ok := new(big.Int).SetString("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 0)
	if !ok {
		panic("unable to parse x coordinate for the default base point (G)")
	}
	y, ok := new(big.Int).SetString("0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 0)
	if !ok {
		panic("unable to parse y coordinate for the default base point (G)")
	}

	return &ec.Point{X: x, Y: y}
}

func createDefaultN() *big.Int {
	n, ok := new(big.Int).SetString("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 0)
	if !ok {
		panic("unable to parse n for the default curve")
	}

	return n
}
