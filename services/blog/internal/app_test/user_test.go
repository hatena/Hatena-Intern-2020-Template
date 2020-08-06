package app

import (
	"context"
	"testing"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_App_Signup_Success(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	sessionExpiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Millisecond).UTC()
	user, sess, err := fixture.app.Signup(ctx, name, password, sessionExpiresAt)
	assert.NoError(t, err)

	user, err = fixture.repo.User().FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, name, user.Name)

	sess, err = fixture.repo.Session().FindByID(ctx, sess.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, sess.UserID)
	assert.Equal(t, sessionExpiresAt, sess.ExpiresAt)
}

func Test_App_Signup_InvalidArgument(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	fixture.accountClient.Error = status.Error(codes.InvalidArgument, "invalid argument")

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	sessionExpiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Millisecond).UTC()
	_, _, err := fixture.app.Signup(ctx, name, password, sessionExpiresAt)
	assert.Equal(t, app.ErrInvalidArgument, err)
}

func Test_App_Signup_AlreadyRegistered(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	fixture.accountClient.Error = status.Error(codes.AlreadyExists, "already registered")

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	sessionExpiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Millisecond).UTC()
	_, _, err := fixture.app.Signup(ctx, name, password, sessionExpiresAt)
	assert.Equal(t, app.ErrAlreadyRegistered, err)
}

func Test_App_Signin_Success(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	sessionExpiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Millisecond).UTC()
	user, sess, err := fixture.app.Signin(ctx, name, password, sessionExpiresAt)
	assert.NoError(t, err)

	user, err = fixture.repo.User().FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, name, user.Name)

	sess, err = fixture.repo.Session().FindByID(ctx, sess.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, sess.UserID)
	assert.Equal(t, sessionExpiresAt, sess.ExpiresAt)
}

func Test_App_Signin_InvalidArgument(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	fixture.accountClient.Error = status.Error(codes.InvalidArgument, "invalid argument")

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	sessionExpiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Millisecond).UTC()
	_, _, err := fixture.app.Signin(ctx, name, password, sessionExpiresAt)
	assert.Equal(t, app.ErrInvalidArgument, err)
}

func Test_App_Signin_AuthenticationFailed(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	fixture.accountClient.Error = status.Error(codes.Unauthenticated, "authentication failed")

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	sessionExpiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Millisecond).UTC()
	_, _, err := fixture.app.Signin(ctx, name, password, sessionExpiresAt)
	assert.Equal(t, app.ErrAuthenticationFailed, err)
}

func Test_App_FindUserBySessionKey_Success(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	sessionExpiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Millisecond).UTC()
	user, sess, err := fixture.app.Signup(ctx, name, password, sessionExpiresAt)
	assert.NoError(t, err)

	u, s, err := fixture.app.FindUserBySessionKey(ctx, sess.Key)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, u.ID)
	assert.Equal(t, sess.ID, s.ID)
}

func Test_App_FindUserBySessionKey_NotFound(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	_, _, err := fixture.app.FindUserBySessionKey(ctx, "unknownkey")
	assert.Error(t, err)
}

func Test_App_FindUserBySessionKey_Expired(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	sessionExpiresAt := time.Now().Add(-24 * time.Hour).Truncate(time.Millisecond).UTC()
	_, sess, err := fixture.app.Signup(ctx, name, password, sessionExpiresAt)
	assert.NoError(t, err)

	_, _, err = fixture.app.FindUserBySessionKey(ctx, sess.Key)
	assert.Error(t, err)
}
