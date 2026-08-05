package kms

import (
	"errors"
	"os"
	"strings"

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

func getMasterKey() ([]byte, error) {
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

func dotenvKey(key string) string {
	key = strings.ReplaceAll(key, ".", "_")
	return strings.ReplaceAll(key, "-", "_")
}

func escapeDotenvValue(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	value = strings.ReplaceAll(value, "\r", `\r`)
	return strings.ReplaceAll(value, "\n", `\n`)
}
