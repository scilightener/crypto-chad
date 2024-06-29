package mathf

import (
	"math/big"
)

func PowerMod(base, exp, mod *big.Int) *big.Int {
	result := big.NewInt(1)
	base = new(big.Int).Mod(base, mod)
	for i := range exp.BitLen() {
		if exp.Bit(i) == 1 {
			result.Mul(result, base)
			result.Mod(result, mod)
		}
		base.Mul(base, base)
		base.Mod(base, mod)
	}

	return result
}

func GCD(a, b *big.Int) (*big.Int, *big.Int, *big.Int) {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	x0, x1, y0, y1 := new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	x1.Set(one)
	y0.Set(one)

	for b.Cmp(zero) != 0 {
		quotient := new(big.Int).Div(a, b)
		a, b = b, new(big.Int).Sub(a, new(big.Int).Mul(quotient, b))
		x0, x1 = x1, new(big.Int).Sub(x0, new(big.Int).Mul(quotient, x1))
		y0, y1 = y1, new(big.Int).Sub(y0, new(big.Int).Mul(quotient, y1))
	}

	return a, x0, y0
}

func ModInv(a, m *big.Int) *big.Int {
	_, _, x := GCD(a, m)
	return x.Mod(x, m)
}
