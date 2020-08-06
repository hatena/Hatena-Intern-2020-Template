package app

import (
	"context"
	"regexp"

	"github.com/hatena/Hatena-Intern-2020/services/account/domain"
	"github.com/hatena/Hatena-Intern-2020/services/account/repository"
)

var userNameRE = regexp.MustCompile(`^[0-9A-Za-z]{3,16}$`)

// Signup は新規ユーザーの登録を行う
func (a *App) Signup(ctx context.Context, name, password string) (*domain.User, error) {
	if ok := userNameRE.MatchString(name); !ok {
		return nil, ErrInvalidArgument
	}
	repo := repository.NewRepository(a.db)
	user, err := domain.CreateUser(name, password)(ctx, repo)
	if err != nil {
		if err == domain.ErrAlreadyExists {
			return nil, ErrAlreadyRegistered
		}
		return nil, err
	}
	return user, nil
}

// Signin はユーザーの認証を行う
func (a *App) Signin(ctx context.Context, name, password string) (*domain.User, error) {
	if ok := userNameRE.MatchString(name); !ok {
		return nil, ErrInvalidArgument
	}
	repo := repository.NewRepository(a.db)
	user, err := repo.User().FindByName(ctx, name)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, ErrAuthenticationFailed
		}
		return nil, err
	}
	success, err := user.Authenticate(password)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, ErrAuthenticationFailed
	}
	return user, nil
}
