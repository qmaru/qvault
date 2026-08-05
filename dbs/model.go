package dbs

type User struct {
	ID         int64  `json:"id" db:"INTEGER PRIMARY KEY AUTOINCREMENT" comment:"user id"`
	Name       string `json:"name" db:"TEXT;DEFAULT ''" comment:"user name"`
	APIKeyHash string `json:"api_key_hash" db:"TEXT;DEFAULT ''" comment:"user API key hash"`
	CreatedAt  int64  `json:"created_at" db:"INTEGER;DEFAULT 0" comment:"user creation timestamp"`
}

type Secret struct {
	ID        int64  `json:"id" db:"INTEGER PRIMARY KEY AUTOINCREMENT" comment:"secret id"`
	UserID    int64  `json:"user_id" db:"INTEGER;DEFAULT 0" comment:"owner user id"`
	Key       string `json:"key" db:"TEXT;DEFAULT ''" comment:"secret key"`
	Value     []byte `json:"value" db:"BLOB;DEFAULT X''" comment:"secret value"`
	CreatedAt int64  `json:"created_at" db:"INTEGER;DEFAULT 0" comment:"secret creation timestamp"`
	UpdatedAt int64  `json:"updated_at" db:"INTEGER;DEFAULT 0" comment:"secret update timestamp"`
}
