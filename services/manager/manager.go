package manager

import (
	"fmt"
	"os"
	"strings"

	"qkms/dbs"
	"qkms/services/common"
	"qkms/services/kms"

	"github.com/joho/godotenv"
	"github.com/qmaru/minitools/v2/secret/chacha20"
)

func ExportToDotenv(apiKey, output string) error {
	if strings.TrimSpace(output) == "" {
		return fmt.Errorf("output path is required")
	}

	userID, err := kms.Authenticate(apiKey)
	if err != nil {
		return err
	}

	masterKey, err := common.GetMasterKey()
	if err != nil {
		return err
	}

	rows, err := dbs.GetDB().Query(fmt.Sprintf(
		"SELECT key, value FROM %s WHERE user_id = ? ORDER BY key", dbs.SecretTable,
	), userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var dotenv strings.Builder
	for rows.Next() {
		var key string
		var encrypted []byte
		if err := rows.Scan(&key, &encrypted); err != nil {
			return err
		}
		key = common.DotenvKey(key)
		key = strings.ToUpper(key)

		value, err := chacha20.New().Decrypt(encrypted, masterKey)
		if err != nil {
			return err
		}
		fmt.Fprintf(&dotenv, "%s=\"%s\"\n", key, common.EscapeDotenvValue(string(value)))
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return os.WriteFile(output, []byte(dotenv.String()), 0600)
}

func ImportFromDotenv(apiKey, input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input path is required")
	}

	userID, err := kms.Authenticate(apiKey)
	if err != nil {
		return err
	}

	env, err := godotenv.Read(input)
	if err != nil {
		return err
	}

	for key, value := range env {
		key = strings.ReplaceAll(key, "_", ".")
		key = strings.ToLower(key)
		if _, err := kms.PutSecret(userID, key, value); err != nil {
			return fmt.Errorf("import key %q: %w", key, err)
		}
	}
	return nil
}

func ListUsers() {}

func ListKeys(apiKey, prefix string) {}
