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

func ExportToDotenv(apiKey, output, prefix string) error {
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
		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
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

func ImportFromDotenv(apiKey, input, prefix string, force bool) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input path is required")
	}

	userID, err := kms.Authenticate(apiKey)
	if err != nil {
		return err
	}

	existingKeys, err := kms.ListKeys(userID)
	if err != nil {
		return err
	}
	existing := make(map[string]struct{}, len(existingKeys))
	for _, key := range existingKeys {
		existing[key] = struct{}{}
	}

	env, err := godotenv.Read(input)
	if err != nil {
		return err
	}

	for key, value := range env {
		key = strings.ReplaceAll(key, "_", ".")
		key = strings.ToLower(key)
		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}
		if !force {
			if _, exists := existing[key]; exists {
				continue
			}
		}
		if _, err := kms.PutSecret(userID, key, value); err != nil {
			return fmt.Errorf("import key %q: %w", key, err)
		}
	}
	return nil
}

func ListKeys(apiKey, prefix string) ([]string, error) {
	userID, err := kms.Authenticate(apiKey)
	if err != nil {
		return nil, err
	}

	keys, err := kms.ListKeys(userID)
	if err != nil {
		return nil, err
	}

	if prefix == "" {
		return keys, nil
	}

	filtered := make([]string, 0, len(keys))
	for _, key := range keys {
		if strings.HasPrefix(key, prefix) {
			filtered = append(filtered, key)
		}
	}
	return filtered, nil
}
