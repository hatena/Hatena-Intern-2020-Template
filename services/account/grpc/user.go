package grpc

import (
	"context"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/account/app"
	"github.com/hatena/Hatena-Intern-2020/services/account/domain"
	pb "github.com/hatena/Hatena-Intern-2020/services/account/pb/account"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) generateToken(user *domain.User) ([]byte, error) {
	claims := jwt.New()
	now := time.Now()
	iss := "hatena-intern-2020-account"
	if err := claims.Set(jwt.IssuerKey, iss); err != nil {
		return nil, err
	}
	sub := "user"
	if err := claims.Set(jwt.SubjectKey, sub); err != nil {
		return nil, err
	}
	exp := now.Add(time.Hour)
	if err := claims.Set(jwt.ExpirationKey, exp); err != nil {
		return nil, err
	}
	if err := claims.Set("user_id", user.ID.String()); err != nil {
		return nil, err
	}
	if err := claims.Set("user_name", user.Name); err != nil {
		return nil, err
	}
	iat := now
	if err := claims.Set(jwt.IssuedAtKey, iat); err != nil {
		return nil, err
	}
	token, err := jwt.Sign(claims, jwa.ES256, s.ecdsaPrivateKey)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Signup は新規ユーザーの登録を行い, トークンを返す
func (s *Server) Signup(ctx context.Context, in *pb.SignupRequest) (*pb.SignupReply, error) {
	user, err := s.app.Signup(ctx, in.Name, in.Password)
	if err != nil {
		if err == app.ErrInvalidArgument {
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}
		if err == app.ErrAlreadyRegistered {
			return nil, status.Error(codes.AlreadyExists, "already registered")
		}
		return nil, err
	}
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}
	return &pb.SignupReply{Token: string(token)}, nil
}

// Signin はユーザーの認証を行い, トークンを返す
func (s *Server) Signin(ctx context.Context, in *pb.SigninRequest) (*pb.SigninReply, error) {
	user, err := s.app.Signin(ctx, in.Name, in.Password)
	if err != nil {
		if err == app.ErrInvalidArgument {
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}
		if err == app.ErrAuthenticationFailed {
			return nil, status.Error(codes.Unauthenticated, "authentication failed")
		}
		return nil, err
	}
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}
	return &pb.SigninReply{Token: string(token)}, nil
}
