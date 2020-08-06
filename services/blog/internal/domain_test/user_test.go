package domain

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_ParseAccountID(t *testing.T) {
	n := testutil.RandUint()
	id, err := domain.ParseAccountID(strconv.FormatUint(uint64(n), 10))
	assert.NoError(t, err)
	assert.Equal(t, domain.AccountID(n), id)
}

func Test_CreateUser_Success(t *testing.T) {
	repo := setup()
	ctx := context.Background()

	accountID := domain.AccountID(testutil.RandUint())
	name := testutil.Rand("user", 8)
	user, err := domain.CreateUser(accountID, name)(ctx, repo)
	assert.NoError(t, err)

	user, err = repo.User().FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, accountID, user.AccountID)
	assert.Equal(t, name, user.Name)
}

func Test_CreateUser_AlreadyExists(t *testing.T) {
	repo := setup()
	ctx := context.Background()

	accountID := domain.AccountID(testutil.RandUint())
	name := testutil.Rand("user", 8)
	_, err := domain.CreateUser(accountID, name)(ctx, repo)
	assert.NoError(t, err)

	_, err = domain.CreateUser(accountID, name)(ctx, repo)
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func Test_User_StartSession_Success(t *testing.T) {
	repo := setup()
	user := testutil.CreateUser(repo)
	ctx := context.Background()

	expiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Millisecond).UTC()
	sess, err := user.StartSession(expiresAt)(ctx, repo)
	assert.NoError(t, err)

	sess, err = repo.Session().FindByID(ctx, sess.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, sess.UserID)
	assert.Equal(t, expiresAt, sess.ExpiresAt)
}

func Test_User_StartSession_KeyRandomness(t *testing.T) {
	repo := setup()
	user := testutil.CreateUser(repo)
	ctx := context.Background()

	expiresAt := time.Now().Add(24 * time.Hour)
	sess1, err := user.StartSession(expiresAt)(ctx, repo)
	assert.NoError(t, err)
	sess2, err := user.StartSession(expiresAt)(ctx, repo)
	assert.True(t, sess1.Key != sess2.Key)
}

func Test_User_CreateBlog_Success(t *testing.T) {
	repo := setup()
	user := testutil.CreateUser(repo)
	ctx := context.Background()

	path := testutil.Rand("blogpath", 8)
	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("blogdesc", 8)
	blog, err := user.CreateBlog(path, title, desc)(ctx, repo)
	assert.NoError(t, err)

	blog, err = repo.Blog().FindByID(ctx, blog.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, blog.UserID)
	assert.Equal(t, path, blog.Path)
	assert.Equal(t, title, blog.Title)
	assert.Equal(t, desc, blog.Description)
}
