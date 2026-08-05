package kms

import (
	"errors"
	"os"

	"qkms/utils"

	"github.com/qmaru/minitools/v2/hashx/blake3"
)

type Secret struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var (
	ErrInvalidAPIKey    = errors.New("invalid api key")
	ErrNotFound         = errors.New("secret not found")
	ErrInvalidMasterKey = errors.New("MASTER_KEY must be exactly 32 bytes")
	ErrInvalidUserName  = errors.New("user name is required")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidSecretKey = errors.New("secret key only allows 0-9, a-z, A-Z, . and -")
)

func masterKey() ([]byte, error) {
	key := os.Getenv("MASTER_KEY")
	return utils.ValidateMasterKey(key)
}

func hashAPIKey(apiKey string) (string, error) {
	hash := blake3.New()
	if _, err := hash.Write([]byte(apiKey)); err != nil {
		return "", err
	}
	return hash.SumStream().ToHex(), nil
}

func ValidateSecretKey(key string) error {
	if key == "" {
		return ErrInvalidSecretKey
	}
	for _, char := range key {
		if (char < '0' || char > '9') &&
			(char < 'a' || char > 'z') &&
			(char < 'A' || char > 'Z') &&
			char != '.' && char != '-' {
			return ErrInvalidSecretKey
		}
	}
	return nil
}
