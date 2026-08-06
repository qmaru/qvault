package kms

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"qkms/dbs"
	"qkms/services/common"
	"qkms/utils"
)

func Authenticate(apiKey string) (int64, error) {
	if apiKey == "" {
		return 0, ErrInvalidAPIKey
	}

	apiKeyHash, err := common.HashAPIKey(apiKey)
	if err != nil {
		return 0, err
	}

	db := dbs.GetDB()

	row, err := db.QueryOne(fmt.Sprintf(
		"SELECT id FROM %s WHERE api_key_hash = ? LIMIT 1", dbs.UserTable,
	), apiKeyHash)
	if err != nil {
		return 0, err
	}

	var userID int64
	if err = row.Scan(&userID); err == sql.ErrNoRows {
		return 0, ErrInvalidAPIKey
	}
	return userID, err
}

func CreateUser(name string, prefix string, rotate bool) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrInvalidUserName
	}

	key, err := utils.GenerateAPIKey()
	if err != nil {
		return "", err
	}
	key = prefix + key

	apiKeyHash, err := common.HashAPIKey(key)
	if err != nil {
		return "", err
	}

	db := dbs.GetDB()
	if rotate {
		result, err := db.Exec(fmt.Sprintf(
			"UPDATE %s SET api_key_hash = ? WHERE name = ?", dbs.UserTable,
		), apiKeyHash, name)
		if err != nil {
			return "", err
		}

		count, err := result.RowsAffected()
		if err != nil {
			return "", err
		}
		if count == 0 {
			return "", ErrUserNotFound
		}
		return key, nil
	}

	_, err = db.Exec(fmt.Sprintf(
		"INSERT INTO %s (name, api_key_hash, created_at) VALUES (?, ?, ?)", dbs.UserTable,
	), name, apiKeyHash, time.Now().Unix())
	if err != nil {
		return "", err
	}
	return key, nil
}
