package ecdh

import (
	"crypto-chad-lib/ec"
	"math/big"
)

type PublicKey struct {
	P *ec.Point
}

func NewPublicKey(p *ec.Point) PublicKey {
	return PublicKey{
		P: ec.NewPoint(p.X, p.Y),
	}
}

type PrivateKey struct {
	K *big.Int
}

func NewPrivateKey(k *big.Int) PrivateKey {
	return PrivateKey{
		K: new(big.Int).Set(k),
	}
}

type Keys struct {
	Public  PublicKey
	Private PrivateKey
}
