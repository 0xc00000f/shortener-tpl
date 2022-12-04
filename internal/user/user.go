package user

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"

	"github.com/0xc00000f/shortener-tpl/internal/crypto"
)

type User struct {
	UserID uuid.UUID
}

var Nil = User{UserID: uuid.Nil}

func New() User {
	return User{UserID: uuid.New()}
}

func (u *User) UserEncrypt() ([]byte, error) {
	aesBlock, err := crypto.NewAESBlock()
	if err != nil {
		return nil, fmt.Errorf("failed user decrypting: %w", err)
	}

	byteUser, err := u.UserID.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed user encrypting: %w", err)
	}

	encryptedUser := make([]byte, aes.BlockSize)
	aesBlock.Encrypt(encryptedUser, byteUser)

	return encryptedUser, nil
}

func (u *User) UserEncryptToString() (string, error) {
	b, err := u.UserEncrypt()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (u *User) UserDecrypt(b []byte) error {
	aesBlock, err := crypto.NewAESBlock()
	if err != nil {
		return fmt.Errorf("failed user decrypting: %w", err)
	}

	dst := make([]byte, len(uuid.UUID{}))
	aesBlock.Decrypt(dst, b)

	userID := uuid.New()
	if err = userID.UnmarshalBinary(dst); err != nil {
		return fmt.Errorf("failed user decrypting: %w", err)
	}

	u.UserID = userID

	return nil
}

func (u *User) UserDecryptFromString(s string) error {
	b, err := hex.DecodeString(s)
	if err != nil {
		return fmt.Errorf("failed user decrypting: %w", err)
	}

	return u.UserDecrypt(b)
}

func Valid(ciphertext string) bool {
	u := User{}
	err := u.UserDecryptFromString(ciphertext)

	return err == nil
}
