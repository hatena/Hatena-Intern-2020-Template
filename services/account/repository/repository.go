package repository

import (
	"github.com/hatena/Hatena-Intern-2020/services/account/domain"
	"github.com/jmoiron/sqlx"
)

// DB はデータベースのインターフェース
type DB interface {
	sqlx.Execer
	sqlx.ExecerContext
	sqlx.Queryer
	sqlx.QueryerContext
}

// Repository は domain.Repository に対するデータベースを使った実装
type Repository struct {
	user *UserRepository
}

// NewRepository は Repository を作成する
func NewRepository(db DB) *Repository {
	return &Repository{
		user: newUserRepository(db),
	}
}

// User はユーザーに対するリポジトリを返す
func (r *Repository) User() domain.UserRepository {
	return r.user
}

func generateID(db DB) (uint64, error) {
	var id uint64
	err := sqlx.Get(db, &id, "SELECT UUID_SHORT()")
	return id, err
}
