package app

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"strings"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	pb_account "github.com/hatena/Hatena-Intern-2020/services/blog/pb/account"
	"github.com/hatena/Hatena-Intern-2020/services/blog/repository"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Signup は新規ユーザーの登録を行う
func (a *App) Signup(ctx context.Context, name, password string, sessionExpiresAt time.Time) (*domain.User, *domain.Session, error) {
	reply, err := a.accountClient.Signup(ctx, &pb_account.SignupRequest{Name: name, Password: password})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, nil, err
		}
		if st.Code() == codes.InvalidArgument {
			return nil, nil, ErrInvalidArgument
		}
		if st.Code() == codes.AlreadyExists {
			return nil, nil, ErrAlreadyRegistered
		}
		return nil, nil, err
	}
	stub, err := verifyToken(reply.Token, a.accountECDSAPublicKey)
	if err != nil {
		return nil, nil, err
	}
	repo := repository.NewRepository(a.db)
	user, err := domain.CreateUser(stub.AccountID, stub.Name)(ctx, repo)
	if err != nil {
		if err == domain.ErrAlreadyExists {
			return nil, nil, ErrAlreadyRegistered
		}
		return nil, nil, err
	}
	sess, err := user.StartSession(sessionExpiresAt)(ctx, repo)
	if err != nil {
		return nil, nil, err
	}
	return user, sess, nil
}

// Signin はユーザーの認証を行う
func (a *App) Signin(ctx context.Context, name, password string, sessionExpiresAt time.Time) (*domain.User, *domain.Session, error) {
	reply, err := a.accountClient.Signin(ctx, &pb_account.SigninRequest{Name: name, Password: password})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, nil, err
		}
		if st.Code() == codes.InvalidArgument {
			return nil, nil, ErrInvalidArgument
		}
		if st.Code() == codes.Unauthenticated {
			return nil, nil, ErrAuthenticationFailed
		}
		return nil, nil, err
	}
	stub, err := verifyToken(reply.Token, a.accountECDSAPublicKey)
	if err != nil {
		return nil, nil, err
	}
	repo := repository.NewRepository(a.db)
	user, err := repo.User().FindByAccountID(ctx, stub.AccountID)
	if err != nil {
		if err != domain.ErrNotFound {
			return nil, nil, err
		}
		user, err = domain.CreateUser(stub.AccountID, stub.Name)(ctx, repo)
		if err != nil {
			if err == domain.ErrAlreadyExists {
				return nil, nil, errors.New("invalid state")
			}
			return nil, nil, err
		}
	}
	sess, err := user.StartSession(sessionExpiresAt)(ctx, repo)
	if err != nil {
		return nil, nil, err
	}

	return user, sess, nil
}

type userStub struct {
	AccountID domain.AccountID
	Name      string
}

func verifyToken(src string, publicKey *ecdsa.PublicKey) (*userStub, error) {
	token, err := jwt.Parse(strings.NewReader(src), jwt.WithVerify(jwa.ES256, publicKey))
	if err != nil {
		return nil, err
	}
	err = jwt.Verify(
		token,
		jwt.WithIssuer("hatena-intern-2020-account"),
		jwt.WithSubject("user"),
		jwt.WithAcceptableSkew(1*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	idRaw, ok := token.Get("user_id")
	if !ok {
		return nil, errors.New("invalid token")
	}
	idStr, ok := idRaw.(string)
	if !ok {
		return nil, errors.New("invalid token")
	}
	id, err := domain.ParseAccountID(idStr)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	nameRaw, ok := token.Get("user_name")
	if !ok {
		return nil, errors.New("invalid token")
	}
	name, ok := nameRaw.(string)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return &userStub{
		AccountID: id,
		Name:      name,
	}, nil
}

// FindUserBySessionKey はセッションのキーからユーザーを取得する
func (a *App) FindUserBySessionKey(ctx context.Context, key string) (*domain.User, *domain.Session, error) {
	repo := repository.NewRepository(a.db)
	sess, err := repo.Session().FindByKey(ctx, key)
	if err != nil {
		return nil, nil, err
	}
	if sess.IsExpired(time.Now()) {
		return nil, nil, errors.New("session expired")
	}
	user, err := repo.User().FindByID(ctx, sess.UserID)
	if err != nil {
		return nil, nil, err
	}
	return user, sess, nil
}
