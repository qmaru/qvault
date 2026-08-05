package kms

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

func ImportFromDotenv(apiKey, input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input path is required")
	}

	userID, err := Authenticate(apiKey)
	if err != nil {
		return err
	}

	env, err := godotenv.Read(input)
	if err != nil {
		return err
	}

	for key, value := range env {
		key = strings.ReplaceAll(key, "_", ".")
		if _, err := PutSecret(userID, key, value); err != nil {
			return fmt.Errorf("import key %q: %w", key, err)
		}
	}
	return nil
}
