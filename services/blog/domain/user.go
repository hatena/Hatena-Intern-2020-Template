package domain

import (
	"context"
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

// UserID はユーザーにユニークに割り当てられる ID
type UserID uint64

// AccountID はアカウントサービス側でユーザーにユニークに割り当てられる ID
type AccountID uint64

// ParseAccountID は文字列の AccountID をパースする
func ParseAccountID(str string) (AccountID, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return AccountID(0), err
	}
	return AccountID(id), nil
}

// User はユーザーを表す
type User struct {
	ID        UserID    `db:"id"`
	AccountID AccountID `db:"account_id"`
	Name      string    `db:"name"`
}

// CreateUserInput はユーザー作成時の入力
type CreateUserInput struct {
	AccountID AccountID
	Name      string
}

// UserRepository はユーザーのリポジトリ
type UserRepository interface {
	Create(ctx context.Context, input *CreateUserInput) (*User, error)
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByAccountID(ctx context.Context, accountID AccountID) (*User, error)
}

// CreateUser は新規ユーザーを作成する
func CreateUser(accountID AccountID, name string) func(ctx context.Context, r Repository) (*User, error) {
	return func(ctx context.Context, r Repository) (*User, error) {
		_, err := r.User().FindByAccountID(ctx, accountID)
		if err != ErrNotFound {
			if err != nil {
				return nil, err
			}
			return nil, ErrAlreadyExists
		}
		return r.User().Create(ctx, &CreateUserInput{
			AccountID: accountID,
			Name:      name,
		})
	}
}

// StartSession は新規セッションを開始する
func (u User) StartSession(expiresAt time.Time) func(ctx context.Context, r Repository) (*Session, error) {
	return func(ctx context.Context, r Repository) (*Session, error) {
		key, err := generateSessionKey(64)
		if err != nil {
			return nil, err
		}
		return r.Session().Create(ctx, &CreateSessionInput{
			UserID:    u.ID,
			Key:       key,
			ExpiresAt: expiresAt,
		})
	}
}

// CreateBlog は新規ブログを作成する
func (u User) CreateBlog(path, title, description string) func(ctx context.Context, r Repository) (*Blog, error) {
	return func(ctx context.Context, r Repository) (*Blog, error) {
		return r.Blog().Create(ctx, &CreateBlogInput{
			UserID:      u.ID,
			Path:        path,
			Title:       title,
			Description: description,
		})
	}
}

var keyLetters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func generateSessionKey(size int) (string, error) {
	seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return "", err
	}
	r := rand.New(rand.NewSource(seed.Int64()))
	key := make([]rune, size)
	for i := range key {
		key[i] = keyLetters[r.Intn(len(keyLetters))]
	}
	return string(key), nil
}
