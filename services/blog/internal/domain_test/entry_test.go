package domain

import (
	"context"
	"testing"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Entry_Edit(t *testing.T) {
	repo := setup()
	user := testutil.CreateUser(repo)
	blog := testutil.CreateBlog(user, repo)
	entry := testutil.CreateEntry(blog, repo)
	ctx := context.Background()

	title := testutil.Rand("entrytitle", 8)
	body := testutil.Rand("entrybody", 8)
	editedAt := time.Now().Add(time.Hour).Truncate(time.Millisecond).UTC()
	renderer := testutil.CreateTestRenderer(func(body string) (string, error) {
		return "<p>" + body + "</p>", nil
	})
	_, err := entry.Edit(title, body, editedAt)(ctx, repo, renderer)
	assert.NoError(t, err)

	edited, err := repo.Entry().FindByID(ctx, entry.ID)
	assert.Equal(t, entry.BlogID, edited.BlogID)
	assert.Equal(t, title, edited.Title)
	assert.Equal(t, body, edited.Body)
	assert.Equal(t, "<p>"+body+"</p>", edited.BodyHTML)
	assert.Equal(t, entry.PublishedAt, edited.PublishedAt)
	assert.Equal(t, editedAt, edited.EditedAt)
}

func Test_Entry_Unpublish(t *testing.T) {
	repo := setup()
	user := testutil.CreateUser(repo)
	blog := testutil.CreateBlog(user, repo)
	entry := testutil.CreateEntry(blog, repo)
	ctx := context.Background()

	err := entry.Unpublish()(ctx, repo)
	assert.NoError(t, err)

	_, err = repo.Entry().FindByID(ctx, entry.ID)
	assert.Equal(t, domain.ErrNotFound, err)
}
