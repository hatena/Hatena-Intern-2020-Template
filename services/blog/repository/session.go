package repository

import (
	"context"
	"database/sql"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/jmoiron/sqlx"
)

// SessionRepository は domain.SessionRepository に対するデータベースを使った実装
type SessionRepository struct {
	db DB
}

func newSessionRepository(db DB) *SessionRepository {
	return &SessionRepository{db}
}

// Create は新規セッションを作成し, リポジトリに保存する
func (r *SessionRepository) Create(ctx context.Context, input *domain.CreateSessionInput) (*domain.Session, error) {
	id, err := generateID(r.db)
	if err != nil {
		return nil, err
	}
	session := &domain.Session{
		ID:        domain.SessionID(id),
		UserID:    input.UserID,
		Key:       input.Key,
		ExpiresAt: input.ExpiresAt,
	}
	_, err = r.db.ExecContext(
		ctx,
		`
			INSERT INTO sessions (id, user_id, `+"`key`"+`, expires_at)
				VALUES (?, ?, ?, ?)
		`,
		session.ID, session.UserID, session.Key, session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// FindByID はリポジトリから ID でセッションを検索する
func (r *SessionRepository) FindByID(ctx context.Context, id domain.SessionID) (*domain.Session, error) {
	var session domain.Session
	err := sqlx.GetContext(
		ctx,
		r.db,
		&session,
		`
			SELECT id, user_id, `+"`key`"+`, expires_at FROM sessions
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
	return &session, nil
}

// FindByKey はリポジトリからキーでセッションを検索する
func (r *SessionRepository) FindByKey(ctx context.Context, key string) (*domain.Session, error) {
	var session domain.Session
	err := sqlx.GetContext(
		ctx,
		r.db,
		&session,
		`
			SELECT id, user_id, `+"`key`"+`, expires_at FROM sessions
				WHERE `+"`key`"+` = ? LIMIT 1
		`,
		key,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &session, nil
}

// Delete はリポジトリからセッションを削除する
func (r *SessionRepository) Delete(ctx context.Context, id domain.SessionID) error {
	_, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM sessions WHERE id = ?
		`,
		id,
	)
	return err
}
