package ec_test

import (
	"crypto-chad-lib/ec"
	"crypto-chad-lib/mathf"
	"fmt"
	"math/big"
	"testing"
)

var testCurve = createTestCurve()

func TestCurve_IsValid(t *testing.T) {
	testCases := []struct {
		name, a, b, p string
		success       bool
	}{
		{"p is not prime", "2", "3", "4", false},
		{"p is not prime, a = 0", "0", "3", "4", false},
		{"p is not prime, b = 0", "2", "0", "4", false},
		{"p is negative", "2", "3", "-1", false},
		{"p = 0", "2", "3", "0", false},
		{"p = 1", "2", "3", "1", false},
		{"p is not prime, a is big", "222222222222222222222222222222222222222", "3", "4", false},
		{"p is not prime, b is big", "2", "333333333333333333333333333333333333333", "4", false},
		{"p is not prime, p is big", "2", "3", "444444444444444444444444444444444444444", false},
		{"(4*n^3 + 27*b^2) % p = 0", "2", "3", "5", false},

		{"p is prime", "2", "3", "7", true},
		{"a = 0", "0", "3", "7", true},
		{"b = 0", "2", "0", "7", true},
		{"p is big", "2", "3", mathf.GeneratePrime(1024).String(), true},
		{"a is big", "222222222222222222222222222222222222222", "3", "7", true},
		{"b is big", "2", "333333333333333333333333333333333333333", "7", true},
		{"a < 0", "-2", "3", "7", true},
		{"b < 0", "2", "-3", "7", true},
		{"test curve", testCurve.A.String(), testCurve.B.String(), testCurve.P.String(), true},
		{"default curve", ec.DefaultCurve.A.String(), ec.DefaultCurve.B.String(), ec.DefaultCurve.P.String(), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			a := parseBigInt(tc.a)
			b := parseBigInt(tc.b)
			p := parseBigInt(tc.p)

			c := ec.NewCurve(a, b, p)
			if c.IsValid() != tc.success {
				t.Fail()
			}
		})
	}
}

func TestCurve_Add(t *testing.T) {
	testCases := []struct {
		name, ax, ay, bx, by, cx, cy string
	}{
		{"a is inf", "", "", "95", "31", "95", "31"},
		{"b is inf", "17", "10", "", "", "17", "10"},
		{"both inf", "", "", "", "", "", ""},
		{"happy path", "17", "10", "95", "31", "1", "54"},
		{"a = b", "17", "10", "17", "10", "32", "90"},
		{"a = -b", "17", "10", "17", "87", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a, b, c *ec.Point
			if len(tc.ax) == 0 && len(tc.ay) == 0 {
				a = ec.NewInfPoint()
			} else {
				x := parseBigInt(tc.ax)
				y := parseBigInt(tc.ay)
				a = ec.NewPoint(x, y)
			}
			if len(tc.bx) == 0 && len(tc.by) == 0 {
				b = ec.NewInfPoint()
			} else {
				x := parseBigInt(tc.bx)
				y := parseBigInt(tc.by)
				b = ec.NewPoint(x, y)
			}
			if len(tc.cx) == 0 && len(tc.cy) == 0 {
				c = ec.NewInfPoint()
			} else {
				x := parseBigInt(tc.cx)
				y := parseBigInt(tc.cy)
				c = ec.NewPoint(x, y)
			}

			curve := testCurve
			res := curve.Add(a, b)
			if !curve.Contains(res) {
				t.Errorf("result not on curve: %v", res)
			}
			if !res.IsEqual(c) {
				t.Errorf("expected %v, got %v", c, res)
			}
		})
	}
}

func TestCurve_Double(t *testing.T) {
	testCases := []struct {
		name, ax, ay, cx, cy string
	}{
		{"a is inf", "", "", "", ""},
		{"happy path", "17", "10", "32", "90"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a, c *ec.Point
			if len(tc.ax) == 0 && len(tc.ay) == 0 {
				a = ec.NewInfPoint()
			} else {
				x := parseBigInt(tc.ax)
				y := parseBigInt(tc.ay)
				a = ec.NewPoint(x, y)
			}
			if len(tc.cx) == 0 && len(tc.cy) == 0 {
				c = ec.NewInfPoint()
			} else {
				x := parseBigInt(tc.cx)
				y := parseBigInt(tc.cy)
				c = ec.NewPoint(x, y)
			}

			curve := testCurve
			res := curve.Double(a)
			if !curve.Contains(res) {
				t.Errorf("result not on curve: %v", res)
			}
			if !res.IsEqual(c) {
				t.Errorf("expected %v, got %v", c, res)
			}
		})
	}
}

func TestCurve_MulScalar(t *testing.T) {
	testCases := []struct {
		name, ax, ay, cx, cy, n string
	}{
		{"a is inf", "", "", "", "", "2"},
		{"happy path", "17", "10", "32", "90", "2"},
		{"n = 0", "17", "10", "", "", "0"},
		{"n = 1", "17", "10", "17", "10", "1"},
		{"n odd", "17", "10", "10", "21", "37"},
		{"n even", "17", "10", "80", "10", "40"},
		{"cy = 0, n odd", "30", "0", "30", "0", "3"},
		{"cy = 0, n even", "30", "0", "", "", "2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a, c *ec.Point
			if len(tc.ax) == 0 && len(tc.ay) == 0 {
				a = ec.NewInfPoint()
			} else {
				x := parseBigInt(tc.ax)
				y := parseBigInt(tc.ay)
				a = ec.NewPoint(x, y)
			}
			if len(tc.cx) == 0 && len(tc.cy) == 0 {
				c = ec.NewInfPoint()
			} else {
				x := parseBigInt(tc.cx)
				y := parseBigInt(tc.cy)
				c = ec.NewPoint(x, y)
			}

			curve := testCurve
			n := parseBigInt(tc.n)
			res := curve.MulScalar(a, n)
			if !curve.Contains(res) {
				t.Errorf("result not on curve: %v", res)
			}
			if !res.IsEqual(c) {
				t.Errorf("expected %v, got %v", c, res)
			}
		})
	}
}

func parseBigInt(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic(fmt.Sprintf("%s is not parsable to big.Int", s))
	}

	return n
}

func createTestCurve() *ec.Curve {
	return ec.NewCurve(big.NewInt(2), big.NewInt(3), big.NewInt(97))
}
