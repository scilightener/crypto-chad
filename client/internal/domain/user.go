package domain

import "crypto-chad-lib/rsa"

type User struct {
	Name string
	Keys *rsa.Keys
}
