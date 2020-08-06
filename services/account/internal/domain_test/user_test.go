package domain

import (
	"context"
	"testing"

	"github.com/hatena/Hatena-Intern-2020/services/account/domain"
	"github.com/hatena/Hatena-Intern-2020/services/account/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_CreateUser_Success(t *testing.T) {
	repo := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	user, err := domain.CreateUser(name, "!passw0rd")(ctx, repo)
	assert.NoError(t, err)
	assert.Equal(t, name, user.Name)
}

func Test_CreateUser_AlreadyExists(t *testing.T) {
	repo := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	_, err := domain.CreateUser(name, "!passw0rd")(ctx, repo)
	assert.NoError(t, err)

	user, err := domain.CreateUser(name, "!passw0rd")(ctx, repo)
	assert.Equal(t, domain.ErrAlreadyExists, err)
	assert.Nil(t, user)
}

func Test_User_Authenticate_Success(t *testing.T) {
	repo := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	user, err := domain.CreateUser(name, "!passw0rd")(ctx, repo)
	assert.NoError(t, err)

	success, err := user.Authenticate("!passw0rd")
	assert.NoError(t, err)
	assert.True(t, success)
}

func Test_User_Authenticate_Failure(t *testing.T) {
	repo := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	user, err := domain.CreateUser(name, "!passw0rd")(ctx, repo)
	assert.NoError(t, err)

	success, err := user.Authenticate("passowrd")
	assert.NoError(t, err)
	assert.False(t, success)
}
