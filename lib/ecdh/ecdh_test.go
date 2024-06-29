package ecdh_test

import (
	"crypto-chad-lib/ecdh"
	"fmt"
	"testing"
)

func TestECDH_ComputeSecret(t *testing.T) {
	const (
		reps = 100
	)
	for i := range reps {
		t.Run(fmt.Sprintf("test compute secret %d", i), func(t *testing.T) {
			t.Parallel()

			alice, err := ecdh.DefaultECDH.GenerateKeys()
			if err != nil {
				t.Fatalf("unable to generate keys: %s", err.Error())
			}

			bob, err := ecdh.DefaultECDH.GenerateKeys()
			if err != nil {
				t.Fatalf("unable to generate keys: %s", err.Error())
			}

			aliceSecret := ecdh.DefaultECDH.ComputeSecret(alice.Private, bob.Public)
			bobSecret := ecdh.DefaultECDH.ComputeSecret(bob.Private, alice.Public)

			if string(aliceSecret) != string(bobSecret) {
				t.Fail()
			}
		})
	}
}
