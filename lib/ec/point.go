package ec

import (
	"fmt"
	"math/big"
)

type Point struct {
	X, Y *big.Int
}

func NewPoint(x, y *big.Int) *Point {
	return &Point{
		X: new(big.Int).Set(x),
		Y: new(big.Int).Set(y),
	}
}

func NewInfPoint() *Point {
	return nil
}

func (p *Point) IsInfinity() bool {
	return p == nil
}

func (p *Point) IsEqual(q *Point) bool {
	if p.IsInfinity() || q.IsInfinity() {
		return p.IsInfinity() && q.IsInfinity()
	}

	return p.X.Cmp(q.X) == 0 && p.Y.Cmp(q.Y) == 0
}

func (p *Point) String() string {
	return fmt.Sprintf("X: %s Y: %s", p.X, p.Y)
}
