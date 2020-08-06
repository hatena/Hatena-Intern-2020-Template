package domain

import (
	"context"
	"testing"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Blog_Edit(t *testing.T) {
	repo := setup()
	user := testutil.CreateUser(repo)
	blog := testutil.CreateBlog(user, repo)
	ctx := context.Background()

	path := blog.Path
	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("blogdesc", 8)
	blog, err := blog.Edit(title, desc)(ctx, repo)
	assert.NoError(t, err)

	blog, err = repo.Blog().FindByID(ctx, blog.ID)
	assert.NoError(t, err)
	assert.Equal(t, path, blog.Path)
	assert.Equal(t, title, blog.Title)
	assert.Equal(t, desc, blog.Description)
}

func Test_Blog_Delete(t *testing.T) {
	repo := setup()
	user := testutil.CreateUser(repo)
	blog := testutil.CreateBlog(user, repo)
	ctx := context.Background()

	err := blog.Delete()(ctx, repo)
	assert.NoError(t, err)

	_, err = repo.Blog().FindByID(ctx, blog.ID)
	assert.Equal(t, domain.ErrNotFound, err)
}

func Test_Blog_PublishEntry(t *testing.T) {
	repo := setup()
	user := testutil.CreateUser(repo)
	blog := testutil.CreateBlog(user, repo)
	ctx := context.Background()

	title := testutil.Rand("entrytitle", 8)
	body := testutil.Rand("entrybody", 8)
	publishedAt := time.Now().Truncate(time.Millisecond).UTC()
	renderer := testutil.CreateTestRenderer(func(body string) (string, error) {
		return "<p>" + body + "</p>", nil
	})
	entry, err := blog.PublishEntry(title, body, publishedAt)(ctx, repo, renderer)
	assert.NoError(t, err)

	entry, err = repo.Entry().FindByID(ctx, entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, blog.ID, entry.BlogID)
	assert.Equal(t, title, entry.Title)
	assert.Equal(t, body, entry.Body)
	assert.Equal(t, "<p>"+body+"</p>", entry.BodyHTML)
	assert.Equal(t, publishedAt, entry.PublishedAt)
	assert.Equal(t, publishedAt, entry.EditedAt)
}
