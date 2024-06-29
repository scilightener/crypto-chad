package ecdh

import (
	"crypto-chad-lib/ec"
	"crypto-chad-lib/hash"
	"crypto/rand"
	"fmt"
	"math/big"
)

type ECDH struct {
	Curve ec.Curve
	G     *ec.Point
	N     *big.Int
}

func NewECDH(curve ec.Curve, g *ec.Point, n *big.Int) *ECDH {
	return &ECDH{
		Curve: curve,
		G:     ec.NewPoint(g.X, g.Y),
		N:     new(big.Int).Set(n),
	}
}

func (ecdh *ECDH) GenerateKeys() (*Keys, error) {
	k, err := rand.Int(rand.Reader, ecdh.N)
	if err != nil {
		return nil, fmt.Errorf("ecdh.GenerateKeys: %s", err.Error())
	}
	p := ecdh.Curve.MulScalar(ecdh.G, k)
	return &Keys{
		Public:  NewPublicKey(p),
		Private: NewPrivateKey(k),
	}, nil
}

func (ecdh *ECDH) ComputeSecret(private PrivateKey, public PublicKey) []byte {
	return hash.Hash(ecdh.Curve.MulScalar(public.P, private.K).String())
}
