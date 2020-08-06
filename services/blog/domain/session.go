package domain

import (
	"context"
	"time"
)

// SessionID はセッションにユニークに割り当てられる ID
type SessionID uint64

// Session はユーザーのセッションを表す
type Session struct {
	ID        SessionID `db:"id"`
	UserID    UserID    `db:"user_id"`
	Key       string    `db:"key"`
	ExpiresAt time.Time `db:"expires_at"`
}

// CreateSessionInput はセッション作成時の入力
type CreateSessionInput struct {
	UserID    UserID
	Key       string
	ExpiresAt time.Time
}

// SessionRepository はセッションのリポジトリ
type SessionRepository interface {
	Create(ctx context.Context, input *CreateSessionInput) (*Session, error)
	FindByID(ctx context.Context, id SessionID) (*Session, error)
	FindByKey(ctx context.Context, key string) (*Session, error)
	Delete(ctx context.Context, id SessionID) error
}

// IsExpired はセッションが有効期限切れかを判定する
func (s Session) IsExpired(now time.Time) bool {
	return s.ExpiresAt.Before(now)
}
