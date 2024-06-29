package mathf_test

import (
	"crypto-chad-lib/mathf"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestPowerMod(t *testing.T) {
	const (
		reps = 1000
		mx   = 1_000_000_000_000_000_000
	)
	for i := range reps {
		t.Run(fmt.Sprintf("test powermod %d", i), func(t *testing.T) {
			t.Parallel()

			base, err := rand.Int(rand.Reader, big.NewInt(mx))
			if err != nil {
				panic(err)
			}

			exp, err := rand.Int(rand.Reader, big.NewInt(mx))
			if err != nil {
				panic(err)
			}

			mod, err := rand.Int(rand.Reader, big.NewInt(mx))
			if err != nil {
				panic(err)
			}

			result := mathf.PowerMod(base, exp, mod)
			expected := new(big.Int).Exp(base, exp, mod)
			if result.Cmp(expected) != 0 {
				t.Errorf("expected %v, got %v", expected, result)
			}
		})
	}
}

func TestGCD(t *testing.T) {
	const (
		reps = 1000
		mx   = 1_000_000_000_000_000_000
	)
	for i := range reps {
		t.Run(fmt.Sprintf("test gcd %d", i), func(t *testing.T) {
			t.Parallel()

			a, err := rand.Int(rand.Reader, big.NewInt(mx))
			if err != nil {
				panic(err)
			}

			b, err := rand.Int(rand.Reader, big.NewInt(mx))
			if err != nil {
				panic(err)
			}

			gcd, s, r := mathf.GCD(a, b)
			x, y := new(big.Int), new(big.Int)
			expected := new(big.Int).GCD(x, y, a, b)
			if gcd.Cmp(expected) != 0 {
				t.Errorf("expected %v, got %v", expected, gcd)
			}
			if r.Cmp(x) != 0 {
				t.Errorf("expected %v, got %v", expected, x)
			}
			if s.Cmp(y) != 0 {
				t.Errorf("expected %v, got %v", expected, y)
			}
		})
	}
}

func TestGeneratePrime(t *testing.T) {
	const (
		reps = 100
		bits = 256
	)
	for i := range reps {
		t.Run(fmt.Sprintf("test generate prime %d", i), func(t *testing.T) {
			t.Parallel()

			prime := mathf.GeneratePrime(bits)
			if !prime.ProbablyPrime(mathf.MillerRabinReps) {
				t.Errorf("expected %v to be prime", prime)
			}
		})
	}
}
