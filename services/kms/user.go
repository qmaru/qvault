package kms

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"qkms/dbs"
	"qkms/utils"
)

func Authenticate(apiKey string) (int64, error) {
	if apiKey == "" {
		return 0, ErrInvalidAPIKey
	}

	apiKeyHash, err := hashAPIKey(apiKey)
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

func CreateUser(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrInvalidUserName
	}

	key, err := utils.GenerateAPIKey()
	if err != nil {
		return "", err
	}
	key = "sk-" + key

	apiKeyHash, err := hashAPIKey(key)
	if err != nil {
		return "", err
	}

	_, err = dbs.GetDB().Exec(fmt.Sprintf(
		"INSERT INTO %s (name, api_key_hash, created_at) VALUES (?, ?, ?)", dbs.UserTable,
	), name, apiKeyHash, time.Now().Unix())
	if err != nil {
		return "", err
	}
	return key, nil
}
