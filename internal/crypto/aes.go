package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func NewAESBlock() (cipher.Block, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	return block, fmt.Errorf("failed creating cipher block: %w", err)
}
