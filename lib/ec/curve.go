package ec

import (
	"crypto-chad-lib/mathf"
	"math/big"
)

// Curve is an elliptic curve of the form: y^2 = x^3 + a*x + b (mod p)
type Curve struct {
	A, B, P *big.Int
}

func NewCurve(a, b, p *big.Int) *Curve {
	return &Curve{
		A: new(big.Int).Set(a),
		B: new(big.Int).Set(b),
		P: new(big.Int).Set(p),
	}
}

func (c Curve) Contains(p *Point) bool {
	if p.IsInfinity() {
		return true
	}

	x3 := mathf.PowerMod(p.X, big.NewInt(3), c.P)
	ax := new(big.Int).Mul(c.A, p.X)
	y := new(big.Int).Mul(p.Y, p.Y)
	rhs := new(big.Int).Add(x3, ax)
	rhs.Add(rhs, c.B)
	rhs.Mod(rhs, c.P)
	y.Mod(y, c.P)
	return y.Cmp(rhs) == 0
}

func (c Curve) IsValid() bool {
	if c.P.Cmp(big.NewInt(0)) <= 0 || !mathf.MillerRabinPrimalityTest(c.P, mathf.MillerRabinReps) {
		return false
	}

	// (4*n^3 + 27*b^2) % p != 0
	a1 := new(big.Int).Mul(big.NewInt(4), mathf.PowerMod(c.A, big.NewInt(3), c.P))
	a2 := new(big.Int).Mul(big.NewInt(27), mathf.PowerMod(c.B, big.NewInt(2), c.P))
	a3 := new(big.Int).Add(a1, a2)
	a3.Mod(a3, c.P)
	return !(a3.IsInt64() && a3.Int64() == 0)
}

func (c Curve) Add(a, b *Point) *Point {
	if a.IsInfinity() {
		return b
	}
	if b.IsInfinity() {
		return a
	}
	if a.X.Cmp(b.X) == 0 {
		if a.Y.Cmp(b.Y) == 0 {
			return c.Double(a)
		}
		return nil
	}

	// lambda = (b.y-n.y) * modInverse(b.x-n.x, c.p) % c.p
	lambda := new(big.Int).Mul(new(big.Int).Sub(b.Y, a.Y), mathf.ModInv(new(big.Int).Sub(b.X, a.X), c.P))
	lambda.Mod(lambda, c.P)

	// x = lambda^2 - n.x - b.x % c.p
	x := new(big.Int).Sub(new(big.Int).Mul(lambda, lambda), a.X)
	x.Sub(x, b.X)
	x.Mod(x, c.P)

	// y = lambda*(n.x-x) - n.y % c.p
	y := new(big.Int).Sub(new(big.Int).Mul(lambda, new(big.Int).Sub(a.X, x)), a.Y)
	y.Mod(y, c.P)
	return &Point{x, y}
}

func (c Curve) Double(a *Point) *Point {
	if a == nil {
		return nil
	}
	if a.Y.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	// lambda = (3*n.x^2 + c.n) * modInverse(2*n.y, c.p) % c.p
	lambda := new(big.Int).Mul(big.NewInt(3), new(big.Int).Mul(a.X, a.X))
	lambda.Add(lambda, c.A)
	lambda.Mul(lambda, mathf.ModInv(new(big.Int).Mul(big.NewInt(2), a.Y), c.P))
	lambda.Mod(lambda, c.P)

	// x = lambda^2 - 2*n.x % c.p
	x := new(big.Int).Sub(new(big.Int).Mul(lambda, lambda), new(big.Int).Mul(big.NewInt(2), a.X))
	x.Mod(x, c.P)

	// y = lambda * (n.x-x) - n.y % c.p
	y := new(big.Int).Sub(new(big.Int).Mul(lambda, new(big.Int).Sub(a.X, x)), a.Y)
	y.Mod(y, c.P)
	return &Point{x, y}
}

func (c Curve) MulScalar(p *Point, a *big.Int) *Point {
	var res *Point
	for i := range a.BitLen() {
		if a.Bit(i) == 1 {
			res = c.Add(res, p)
		}
		p = c.Double(p)
	}
	return res
}
