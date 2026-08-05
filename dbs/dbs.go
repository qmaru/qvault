package dbs

import (
	"fmt"
	"os"
	"sync"

	"github.com/qmaru/qdb/sqlitep"
)

const (
	UserTable   = "user"
	SecretTable = "secret"
)

var GetDB = sync.OnceValue(func() *sqlitep.Sqlite {
	return sqlitep.New(os.Getenv("SQLITE_DB_PATH"))
})

func CreateIndexes() error {
	db := GetDB()

	return db.Transaction(func(tx sqlitep.Tx) error {
		_, err := db.ExecWithTx(tx, "DROP INDEX IF EXISTS idx_user_api_key_hash;")
		if err != nil {
			return err
		}
		_, err = db.ExecWithTx(tx, "DROP INDEX IF EXISTS idx_secret_user_key;")
		if err != nil {
			return err
		}

		userIndex := fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS idx_user_api_key_hash ON %s(api_key_hash);", UserTable)
		_, err = db.ExecWithTx(tx, userIndex)
		if err != nil {
			return err
		}

		secretIndex := fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS idx_secret_user_key ON %s(user_id,key);", SecretTable)
		_, err = db.ExecWithTx(tx, secretIndex)
		return err
	})
}

func CreateDB() error {
	db := GetDB()

	tables := []any{
		User{},
		Secret{},
	}

	return db.CreateTable(tables)
}
