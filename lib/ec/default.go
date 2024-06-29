package ec

import "math/big"

// DefaultCurve is secp256k1 curve
var (
	DefaultCurve = createDefaultCurve()
)

func createDefaultCurve() *Curve {
	p, ok := new(big.Int).SetString("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f", 0)
	if !ok {
		panic("unable to parse prime for the default curve")
	}
	a := big.NewInt(0)
	b := big.NewInt(7)
	return NewCurve(a, b, p)
}
