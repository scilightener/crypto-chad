package mathf

import (
	"crypto/rand"
	"math/big"
)

const (
	MillerRabinReps = 20
)

func MillerRabinPrimalityTest(n *big.Int, reps int) bool {
	if n.Cmp(big.NewInt(1)) <= 0 {
		return false
	}
	if n.Cmp(big.NewInt(3)) <= 0 {
		return true
	}
	one := big.NewInt(1)
	two := big.NewInt(2)
	n1 := new(big.Int).Sub(n, one)
	n3 := new(big.Int).Sub(n, big.NewInt(3))
	for i := 0; i < reps; i++ {
		a, err := rand.Int(rand.Reader, n3)
		if err != nil {
			panic(err)
		}
		a.Add(a, two)
		p := PowerMod(a, n1, n)
		if p.Cmp(one) != 0 {
			return false
		}
	}
	return true
}

func GeneratePrime(bits int) *big.Int {
	if bits < 2 {
		panic("mathf.prime.GeneratePrime: bits should be >= 2")
	}
	base := new(big.Int).Lsh(big.NewInt(1), uint(bits))
	mx := new(big.Int).Lsh(big.NewInt(1), uint(bits-1))
	p, err := rand.Int(rand.Reader, mx)
	if err != nil {
		panic(err)
	}
	p.Add(p, base)
	two := big.NewInt(2)
	if new(big.Int).Mod(p, two).Cmp(big.NewInt(1)) != 0 {
		p.Add(p, big.NewInt(1))
	}
	for {
		if MillerRabinPrimalityTest(p, MillerRabinReps) {
			return p
		}
		p.Add(p, two)
	}
}
