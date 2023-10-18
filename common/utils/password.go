package utils

import (
	"crypto/sha512"
	"github.com/anaskhan96/go-password-encoder"
)

var (
	Options *password.Options
)

func init() {
	Options = &password.Options{16, 100, 32, sha512.New}
}

func GeneratePassword(code string) (string, string) {
	return password.Encode(code, Options)
}

func CheckPassword(code string, salt string, encode string) bool {
	return password.Verify(code, salt, encode, Options)
}
