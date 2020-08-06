package repository

import (
	"context"
	"database/sql"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/jmoiron/sqlx"
)

// UserRepository は domain.UserRepository に対するデータベースを使った実装
type UserRepository struct {
	db DB
}

func newUserRepository(db DB) *UserRepository {
	return &UserRepository{db}
}

// Create は新規ユーザーを作成し, リポジトリに保存する
func (r *UserRepository) Create(ctx context.Context, input *domain.CreateUserInput) (*domain.User, error) {
	id, err := generateID(r.db)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		ID:        domain.UserID(id),
		AccountID: input.AccountID,
		Name:      input.Name,
	}
	_, err = r.db.ExecContext(
		ctx,
		`
			INSERT INTO users (id, account_id, name)
				VALUES (?, ?, ?)
		`,
		user.ID, user.AccountID, user.Name,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByID はリポジトリから ID でユーザーを検索する
func (r *UserRepository) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var user domain.User
	err := sqlx.GetContext(
		ctx,
		r.db,
		&user,
		`
			SELECT id, account_id, name FROM users
				WHERE id = ? LIMIT 1
		`,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByAccountID はリポジトリからアカウントサービスの ID でユーザーを検索する
func (r *UserRepository) FindByAccountID(ctx context.Context, id domain.AccountID) (*domain.User, error) {
	var user domain.User
	err := sqlx.GetContext(
		ctx,
		r.db,
		&user,
		`
			SELECT id, account_id, name FROM users
				WHERE account_id = ? LIMIT 1
		`,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
