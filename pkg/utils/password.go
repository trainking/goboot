package utils

import (
	"crypto/sha256"
	"fmt"

	"github.com/trainking/goboot/pkg/random"
)

// Sha256PasswordEncrypt 对密码进行sha256加密
func Sha256PasswordEncrypt(password string) (string, string, error) {
	salt, err := random.RandStringHex(8)
	if err != nil {
		return "", "", err
	}

	h := sha256.New()
	h.Write([]byte(password + salt))

	return fmt.Sprintf("%x", h.Sum(nil)), salt, nil
}

// Sha256PasswordValidate 对密码进行sha256加密，并与已加密的密码进行比较
func Sha256PasswordValidate(password, salt, hash string) bool {
	h := sha256.New()
	h.Write([]byte(password + salt))

	return fmt.Sprintf("%x", h.Sum(nil)) == hash
}
