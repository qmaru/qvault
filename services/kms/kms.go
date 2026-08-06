package kms

import (
	"errors"
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
