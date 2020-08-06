package app

import (
	"github.com/jmoiron/sqlx"
)

// App はアプリケーションを表す
type App struct {
	db *sqlx.DB
}

// NewApp は App を作成する
func NewApp(db *sqlx.DB) *App {
	return &App{db}
}
