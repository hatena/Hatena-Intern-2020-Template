package grpc

import (
	"context"
	"testing"

	"github.com/hatena/Hatena-Intern-2020/services/account/internal/testutil"
	pb "github.com/hatena/Hatena-Intern-2020/services/account/pb/account"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Server_Signup_Success(t *testing.T) {
	s, app := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	reply, err := s.Signup(ctx, &pb.SignupRequest{
		Name:     name,
		Password: "!passw0rd",
	})
	assert.NoError(t, err)

	token, err := verifyToken(reply.Token)
	assert.NoError(t, err)
	assert.Equal(t, "hatena-intern-2020-account", token.Issuer())
	assert.Equal(t, "user", token.Subject())

	userID, ok := token.Get("user_id")
	assert.True(t, ok)
	user, err := app.Signin(ctx, name, "!passw0rd")
	assert.NoError(t, err)
	assert.Equal(t, user.ID.String(), userID)

	n, ok := token.Get("user_name")
	assert.True(t, ok)
	assert.Equal(t, name, n)
}

func Test_Server_Signup_AlreadyRegistered(t *testing.T) {
	s, app := setup()
	ctx := context.Background()
	name := testutil.Rand("user", 8)
	_, err := app.Signup(ctx, name, "!passw0rd")
	assert.NoError(t, err)

	reply, err := s.Signup(ctx, &pb.SignupRequest{
		Name:     name,
		Password: "!passw0rd",
	})
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Nil(t, reply)
}

func Test_Server_Signin_Success(t *testing.T) {
	s, app := setup()
	ctx := context.Background()
	name := testutil.Rand("user", 8)
	user, err := app.Signup(ctx, name, "!passw0rd")
	assert.NoError(t, err)

	reply, err := s.Signin(ctx, &pb.SigninRequest{
		Name:     name,
		Password: "!passw0rd",
	})
	assert.NoError(t, err)

	token, err := verifyToken(reply.Token)
	assert.NoError(t, err)
	assert.Equal(t, "hatena-intern-2020-account", token.Issuer())
	assert.Equal(t, "user", token.Subject())

	userID, ok := token.Get("user_id")
	assert.True(t, ok)
	assert.Equal(t, user.ID.String(), userID)

	n, ok := token.Get("user_name")
	assert.True(t, ok)
	assert.Equal(t, name, n)
}

func Test_Server_Signin_Failure(t *testing.T) {
	s, app := setup()
	ctx := context.Background()
	name := testutil.Rand("user", 8)
	_, err := app.Signup(ctx, name, "!passw0rd")
	assert.NoError(t, err)

	reply, err := s.Signin(ctx, &pb.SigninRequest{
		Name:     name,
		Password: "password",
	})
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Nil(t, reply)
}
