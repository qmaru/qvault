package common

import (
	"os"
	"strings"

	"qkms/utils"

	"github.com/qmaru/minitools/v2/hashx/blake3"
)

func GetMasterKey() ([]byte, error) {
	key := os.Getenv("MASTER_KEY")
	return utils.ValidateMasterKey(key)
}

func HashAPIKey(apiKey string) (string, error) {
	hash := blake3.New()
	if _, err := hash.Write([]byte(apiKey)); err != nil {
		return "", err
	}
	return hash.SumStream().ToHex(), nil
}

func DotenvKey(key string) string {
	key = strings.ReplaceAll(key, ".", "_")
	return strings.ReplaceAll(key, "-", "_")
}

func EscapeDotenvValue(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	value = strings.ReplaceAll(value, "\r", `\r`)
	return strings.ReplaceAll(value, "\n", `\n`)
}
