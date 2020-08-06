package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Connect は DB ハンドルを作成し, 接続が確立するのを待つ
func Connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	attempts := 0
	maxAttempts := 10
	delay := time.Millisecond * 500
	for {
		attempts++

		var got int
		err = db.Get(&got, `select 1`)
		if err == nil {
			break
		}

		if maxAttempts <= attempts {
			return nil, fmt.Errorf("%+v (after %d attempts)", err, attempts)
		}

		time.Sleep(delay)
		delay = time.Duration(1.5 * float64(delay))
	}

	return db, nil
}
