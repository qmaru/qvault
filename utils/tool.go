package utils

import (
	"fmt"
	"os"

	"github.com/qmaru/minitools/v2/encoding/text"
	"github.com/qmaru/minitools/v2/random/nanoid"

	"github.com/joho/godotenv"
)

const ALPHABET string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func GenerateAPIKey() (string, error) {
	n := nanoid.New()
	apiKey, err := n.Generate(ALPHABET, 43)
	if err != nil {
		return "", err
	}
	return apiKey, nil
}

func GenerateMasterKey() (string, error) {
	t := text.New()
	nonce, err := t.Nonce(32)
	if err != nil {
		return "", err
	}

	return t.HexEncode(nonce), nil
}

func ValidateMasterKey(key string) ([]byte, error) {
	t := text.New()

	h, err := t.HexDecoding(key)

	if err != nil {
		return nil, err
	}

	if len(h.DecodeString()) != 32 {
		return nil, fmt.Errorf("invalid master key length: %d", len(h.DecodeString()))
	}
	return h.DecodeRaw(), nil
}
