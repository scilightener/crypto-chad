package rsa

import (
	"crypto-chad-lib/mathf"
	"math/big"
)

const (
	PrimeNumberBits = 1024
)

type PublicKey struct {
	E, N *big.Int
}

type PrivateKey struct {
	D, N *big.Int
}

type Keys struct {
	PublicKey  *PublicKey
	PrivateKey *PrivateKey
}

func NewKeys() *Keys {
	p, q := mathf.GeneratePrime(PrimeNumberBits), mathf.GeneratePrime(PrimeNumberBits)
	n := new(big.Int).Mul(p, q)

	one := big.NewInt(1)
	phi := new(big.Int).Mul(p.Sub(p, one), q.Sub(q, one))
	e, d := generateParams(phi)

	pub := &PublicKey{E: e, N: n}
	priv := &PrivateKey{D: d, N: n}

	return &Keys{
		PublicKey:  pub,
		PrivateKey: priv,
	}
}

func generateParams(phi *big.Int) (*big.Int, *big.Int) {
	e := big.NewInt(3)
	g, _, d := mathf.GCD(e, phi)
	one := big.NewInt(1)
	two := big.NewInt(2)
	for g.Cmp(one) != 0 {
		e.Add(e, two)
		g, _, d = mathf.GCD(e, phi)
	}

	return e, d.Mod(d, phi)
}
