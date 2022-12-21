package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func NewAESBlock() (cipher.Block, error) {
	return aes.NewCipher([]byte(secretKey))
}
