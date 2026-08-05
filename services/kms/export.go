package kms

import (
	"fmt"
	"os"
	"strings"

	"qkms/dbs"

	"github.com/qmaru/minitools/v2/secret/chacha20"
)

func ExportToDotenv(apiKey, output string) error {
	if strings.TrimSpace(output) == "" {
		return fmt.Errorf("output path is required")
	}

	userID, err := Authenticate(apiKey)
	if err != nil {
		return err
	}

	masterKey, err := getMasterKey()
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
		key = dotenvKey(key)
		key = strings.ToUpper(key)

		value, err := chacha20.New().Decrypt(encrypted, masterKey)
		if err != nil {
			return err
		}
		fmt.Fprintf(&dotenv, "%s=\"%s\"\n", key, escapeDotenvValue(string(value)))
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return os.WriteFile(output, []byte(dotenv.String()), 0600)
}
