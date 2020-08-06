package app

import (
	"context"
	"testing"

	"github.com/hatena/Hatena-Intern-2020/services/account/app"
	"github.com/hatena/Hatena-Intern-2020/services/account/domain"
	"github.com/hatena/Hatena-Intern-2020/services/account/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_App_Signup_Success(t *testing.T) {
	a, repo := setup()
	ctx := context.Background()

	name := testutil.Rand("user", 8)
	user, err := a.Signup(ctx, name, "!passw0rd")
	assert.NoError(t, err)

	_, err = repo.User().FindByID(ctx, user.ID)
	assert.NoError(t, err)
}

func Test_App_Signup_AlreadyRegistered(t *testing.T) {
	a, repo := setup()
	ctx := context.Background()
	name := testutil.Rand("user", 8)
	domain.CreateUser(name, "!passw0rd")(ctx, repo)

	_, err := a.Signup(ctx, name, "!passw0rd")
	assert.Equal(t, app.ErrAlreadyRegistered, err)
}

func Test_App_Signin_Success(t *testing.T) {
	a, repo := setup()
	ctx := context.Background()
	name := testutil.Rand("user", 8)
	domain.CreateUser(name, "!passw0rd")(ctx, repo)

	user, err := a.Signin(ctx, name, "!passw0rd")
	assert.NoError(t, err)
	assert.Equal(t, name, user.Name)
}

func Test_App_Signin_Failure(t *testing.T) {
	a, repo := setup()
	ctx := context.Background()
	name := testutil.Rand("user", 8)
	domain.CreateUser(name, "!passw0rd")(ctx, repo)

	_, err := a.Signin(ctx, name, "password")
	assert.Equal(t, app.ErrAuthenticationFailed, err)
}
