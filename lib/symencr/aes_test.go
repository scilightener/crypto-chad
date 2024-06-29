package symencr_test

import (
	"crypto-chad-lib/rnd"
	"crypto-chad-lib/symencr"
	"fmt"
	"testing"
)

const (
	smallKeySize = 8 * (iota + 2)
	mediumKeySize
	largeKeySize
)

func TestAES_Functional(t *testing.T) {
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

	key := []byte(rnd.String(mediumKeySize))
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			encrypted := symencr.Encrypt([]byte(tc.message), key)
			decrypted := symencr.Decrypt(encrypted, key)
			if string(decrypted) != tc.message {
				t.Fail()
			}
		})
	}
}

func TestAES_Functional_SmallKeySize(t *testing.T) {
	const (
		reps        = 1000
		messageSize = 1000
	)

	for i := range reps {
		t.Run(fmt.Sprintf("aes functional test %d small key", i), func(t *testing.T) {
			t.Parallel()

			key := []byte(rnd.String(smallKeySize))
			message := []byte(rnd.String(messageSize))

			encrypted := symencr.Encrypt(message, key)
			decrypted := symencr.Decrypt(encrypted, key)
			if string(decrypted) != string(message) {
				t.Fail()
			}
		})
	}
}

func TestAES_Functional_MediumKeySize(t *testing.T) {
	const (
		reps        = 1000
		messageSize = 1000
	)

	for i := range reps {
		t.Run(fmt.Sprintf("aes functional test %d medium key", i), func(t *testing.T) {
			t.Parallel()

			key := []byte(rnd.String(mediumKeySize))
			message := []byte(rnd.String(messageSize))

			encrypted := symencr.Encrypt(message, key)
			decrypted := symencr.Decrypt(encrypted, key)
			if string(decrypted) != string(message) {
				t.Fail()
			}
		})
	}
}

func TestAES_Functional_LargeKeySize(t *testing.T) {
	const (
		reps        = 1000
		messageSize = 1000
	)

	for i := range reps {
		t.Run(fmt.Sprintf("aes functional test %d large key", i), func(t *testing.T) {
			t.Parallel()

			key := []byte(rnd.String(largeKeySize))
			message := []byte(rnd.String(messageSize))

			encrypted := symencr.Encrypt(message, key)
			decrypted := symencr.Decrypt(encrypted, key)
			if string(decrypted) != string(message) {
				t.Fail()
			}
		})
	}
}

func TestEncrypt_ChangesMessage(t *testing.T) {
	const (
		messageSize = 1000
	)

	message := []byte(rnd.String(messageSize))
	key := []byte(rnd.String(largeKeySize))

	encrypted := symencr.Encrypt(message, key)
	if string(encrypted) == string(message) {
		t.Error("expected encrypted message to be different from original message")
	}
}
