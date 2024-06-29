package rsa_test

import (
	"crypto-chad-lib/rnd"
	"crypto-chad-lib/rsa"
	"fmt"
	"testing"
)

var (
	testKeys = rsa.NewKeys()
)

func TestRSA_Functional(t *testing.T) {
	testCases := []struct {
		name, message string
	}{
		{"small message len=0", ""},
		{"small message len=1", "1"},
		{"avg message len=50", "123456789 123456789 123456789 123456789 123456789 "},
		{"avg message len=128", "1234567 _234567 __34567 ___4567 ____567 _____67 ______7 _______ " +
			"1234567 _234567 __34567 ___4567 ____567 _____67 ______7 _______ "},
		{"non-ascii symbols russian", "сообщение на русском"},
		{"non-ascii symbols chinese", "消息 的英文"},
		{"non-ascii symbols arabic", "رسالة باللغة العربية"},
		{"non-ascii symbols emojis", "😊😎"},
		{"non-ascii symbols kaomojis", "(❁´◡`❁)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			encrypted := rsa.Encrypt([]byte(tc.message), testKeys.PrivateKey.D, testKeys.PrivateKey.N)
			decrypted := rsa.Decrypt(encrypted, testKeys.PublicKey.E, testKeys.PublicKey.N)
			if string(decrypted) != tc.message {
				t.Fail()
			}
		})
	}
}

func TestRSA_FunctionalBigMessages(t *testing.T) {
	const (
		repeats       = 500
		messageLength = 1000
	)
	for i := range repeats {
		message := rnd.String(messageLength)
		t.Run(fmt.Sprintf("functional test big message repeat %d", i), func(t *testing.T) {
			t.Parallel()

			encrypted := rsa.Encrypt([]byte(message), testKeys.PrivateKey.D, testKeys.PrivateKey.N)
			decrypted := rsa.Decrypt(encrypted, testKeys.PublicKey.E, testKeys.PublicKey.N)
			if string(decrypted) != message {
				t.Fail()
			}
		})
	}
}

func TestEncrypt_ChangesMessage(t *testing.T) {
	const (
		messageSize = 1000
	)

	message := rnd.String(messageSize)
	encrypted := rsa.Encrypt([]byte(message), testKeys.PrivateKey.D, testKeys.PrivateKey.N)
	if string(mergeSlices(encrypted)) == message {
		t.Error("expected encrypted message to be different from original message")
	}
}

func mergeSlices(slices [][]byte) []byte {
	var result []byte
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func TestRSA_Chunks(t *testing.T) {
	testCases := []struct {
		name, message string
		nChunks       int
	}{
		{"0 chunk", "", 0},
		{"1 chunk len=128", "hello", 1},

		{"1 chunk len=128", "1234567 _234567 __34567 ___4567 ____567 _____67 ______7 _______ " +
			"1234567 _234567 __34567 ___4567 ____567 _____67 ______7 _______ ", 1},

		{"2 blocks len=129", "+1234567 _234567 __34567 ___4567 ____567 _____67 ______7 _______ " +
			"1234567 _234567 __34567 ___4567 ____567 _____67 ______7 _______ ", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := []byte(tc.message)
			encrypted := rsa.Encrypt(m, testKeys.PrivateKey.D, testKeys.PrivateKey.N)
			if len(encrypted) != tc.nChunks {
				t.Errorf("expected %d chunks, got %d", tc.nChunks, len(encrypted))
			}
		})
	}
}
