package secret

import (
	"database/sql"
	"fmt"
	"time"

	"qvault/dbs"
	"qvault/services/common"

	"github.com/qmaru/minitools/v2/secret/chacha20"
	"github.com/qmaru/qdb/sqlitep"
)

func ListKeys(userID int64) ([]string, error) {
	db := dbs.GetDB()

	rows, err := db.Query(fmt.Sprintf(
		"SELECT key FROM %s WHERE user_id = ? ORDER BY key", dbs.SecretTable,
	), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]string, 0)
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func GetSecret(userID int64, key string) (*Secret, error) {
	if err := ValidateSecretKey(key); err != nil {
		return nil, err
	}

	masterKey, err := common.GetMasterKey()
	if err != nil {
		return nil, err
	}

	db := dbs.GetDB()

	row, err := db.QueryOne(fmt.Sprintf(
		"SELECT value FROM %s WHERE user_id = ? AND key = ? LIMIT 1", dbs.SecretTable,
	), userID, key)
	if err != nil {
		return nil, err
	}

	var encrypted []byte
	if err = row.Scan(&encrypted); err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	plain, err := chacha20.New().Decrypt(encrypted, masterKey)
	if err != nil {
		return nil, err
	}
	return &Secret{Key: key, Value: string(plain)}, nil
}

func PutSecret(userID int64, key, value string) (*Secret, error) {
	if err := ValidateSecretKey(key); err != nil {
		return nil, err
	}

	masterKey, err := common.GetMasterKey()
	if err != nil {
		return nil, err
	}

	nonce, err := chacha20.New().GenerateNonce()
	if err != nil {
		return nil, err
	}
	encrypted, err := chacha20.New().Encrypt([]byte(value), masterKey, nonce)
	if err != nil {
		return nil, err
	}

	db := dbs.GetDB()

	now := time.Now().Unix()
	err = db.Transaction(func(tx sqlitep.Tx) error {
		_, err := db.ExecWithTx(tx, fmt.Sprintf(
			"INSERT INTO %s (user_id, key, value, created_at, updated_at) VALUES (?, ?, ?, ?, ?) "+
				"ON CONFLICT(user_id, key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at",
			dbs.SecretTable,
		), userID, key, encrypted, now, now)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Secret{Key: key, Value: value}, nil
}

func DeleteSecret(userID int64, key string) error {
	if err := ValidateSecretKey(key); err != nil {
		return err
	}

	db := dbs.GetDB()

	result, err := db.Exec(fmt.Sprintf(
		"DELETE FROM %s WHERE user_id = ? AND key = ?", dbs.SecretTable,
	), userID, key)
	if err != nil {
		return err
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return ErrNotFound
	}
	return nil
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
