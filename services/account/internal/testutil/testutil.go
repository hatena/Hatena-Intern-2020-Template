package testutil

import (
	"errors"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL ドライバを使うために必要
	"github.com/hatena/Hatena-Intern-2020/services/account/db"
	"github.com/jmoiron/sqlx"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewDB はテスト用に DB ハンドルを作成する
func NewDB() (*sqlx.DB, error) {
	databaseDSN := os.Getenv("TEST_DATABASE_DSN")
	if databaseDSN == "" {
		return nil, errors.New("TEST_DATABASE_DSN is not set")
	}
	db, err := db.Connect(databaseDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}

var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// Rand はランダムな文字列を返す
func Rand(prefix string, size int) string {
	key := make([]rune, size)
	for i := range key {
		key[i] = letters[rand.Intn(len(letters))]
	}
	return prefix + string(key)
}
