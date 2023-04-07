package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	// 默认长度
	DefaultAesBlockSize = 16
)

// AES ECB模式的加密解密
type AesTool struct {
	//128 192  256位的其中一个 长度 对应分别是 16 24  32字节长度
	Key       []byte
	BlockSize int
}

func NewAesTool(key []byte, blockSize int) *AesTool {
	return &AesTool{Key: key, BlockSize: blockSize}
}

// EncryptBase64 加密并返回base64编码
func (a *AesTool) Encrypt(plaintext []byte) ([]byte, error) {

	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, a.BlockSize+len(plaintext))
	iv := ciphertext[:a.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[a.BlockSize:], plaintext)

	// convert to base64
	return []byte(base64.URLEncoding.EncodeToString(ciphertext)), nil
}

// Decrypt 从base64编码中返回
func (a *AesTool) Decrypt(src []byte) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(string(src))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < a.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:a.BlockSize]
	ciphertext = ciphertext[a.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}
