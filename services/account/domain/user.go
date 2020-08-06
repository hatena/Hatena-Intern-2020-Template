package domain

import (
	"context"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// UserID はユーザーにユニークに割り当てられる ID
type UserID uint64

func (id UserID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

// User はユーザーを表す
type User struct {
	ID           UserID `db:"id"`
	Name         string `db:"name"`
	PasswordHash string `db:"password_hash"`
}

// CreateUserInput はユーザー作成時の入力
type CreateUserInput struct {
	Name         string
	PasswordHash string
}

// UserRepository はユーザーのリポジトリ
type UserRepository interface {
	Create(ctx context.Context, input *CreateUserInput) (*User, error)
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
}

// CreateUser は新規ユーザーを作成する
func CreateUser(name, password string) func(ctx context.Context, r Repository) (*User, error) {
	return func(ctx context.Context, r Repository) (*User, error) {
		_, err := r.User().FindByName(ctx, name)
		if err != ErrNotFound {
			if err != nil {
				return nil, err
			}
			return nil, ErrAlreadyExists
		}
		passwordHash, err := calcPasswordHash(password)
		if err != nil {
			return nil, err
		}
		return r.User().Create(ctx, &CreateUserInput{
			Name:         name,
			PasswordHash: passwordHash,
		})
	}
}

// Authenticate はパスワードによるユーザー認証を行う
func (u User) Authenticate(password string) (bool, error) {
	return comparePasswordHash(u.PasswordHash, password)
}

func calcPasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func comparePasswordHash(passwordHash, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
