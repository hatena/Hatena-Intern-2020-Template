package app

import (
	"context"
	"testing"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_App_ListEntriesByBlog_All(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	blog2 := testutil.CreateBlog(user, fixture.repo)
	entry1 := testutil.CreateEntry(blog, fixture.repo)
	entry2 := testutil.CreateEntry(blog, fixture.repo)
	entry3 := testutil.CreateEntry(blog, fixture.repo)
	entry4 := testutil.CreateEntry(blog2, fixture.repo)
	ctx := context.Background()

	entries, hasNextPage, err := fixture.app.ListEntriesByBlog(ctx, blog, 1, 5)
	assert.NoError(t, err)
	assert.Len(t, entries, 3)
	assert.Contains(t, entries, entry1)
	assert.Contains(t, entries, entry2)
	assert.Contains(t, entries, entry3)
	assert.NotContains(t, entries, entry4)
	assert.False(t, hasNextPage)
}

func Test_App_ListEntriesByBlog_Paging(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	testutil.CreateEntry(blog, fixture.repo)
	testutil.CreateEntry(blog, fixture.repo)
	testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	entries, hasNextPage, err := fixture.app.ListEntriesByBlog(ctx, blog, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.True(t, hasNextPage)

	entries, hasNextPage, err = fixture.app.ListEntriesByBlog(ctx, blog, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.False(t, hasNextPage)
}

func Test_App_PublishEntry_Success(t *testing.T) {
	fixture := setup()
	fixture.rendererClient.RenderFunc = func(src string) string {
		return "<p>" + src + "</p>"
	}
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("entrytitle", 8)
	body := testutil.Rand("entrybody", 8)
	entry, err := fixture.app.PublishEntry(ctx, user, blog, title, body)
	assert.NoError(t, err)

	entry, err = fixture.repo.Entry().FindByID(ctx, entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, blog.ID, entry.BlogID)
	assert.Equal(t, title, entry.Title)
	assert.Equal(t, body, entry.Body)
	assert.Equal(t, "<p>"+body+"</p>", entry.BodyHTML)
}

func Test_App_PublishEntry_NoPermission(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	user2 := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("entrytitle", 8)
	body := testutil.Rand("entrybody", 8)
	_, err := fixture.app.PublishEntry(ctx, user2, blog, title, body)
	assert.Equal(t, app.ErrPermissionDenied, err)
}

func Test_App_PublishEntry_TooLongTitle(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("", 501)
	body := testutil.Rand("entrybody", 8)
	_, err := fixture.app.PublishEntry(ctx, user, blog, title, body)
	assert.Equal(t, app.ErrInvalidArgument, err)
}

func Test_App_FindEntryByID_Found(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	found, err := fixture.app.FindEntryByID(ctx, blog, entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, entry.ID, found.ID)
}

func Test_App_FindEntryByID_NotFoundInBlog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	blog2 := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	_, err := fixture.app.FindEntryByID(ctx, blog2, entry.ID)
	assert.Equal(t, app.ErrNotFound, err)
}

func Test_App_EditEntry_Success(t *testing.T) {
	fixture := setup()
	fixture.rendererClient.RenderFunc = func(src string) string {
		return "<p>" + src + "</p>"
	}
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("entrytitle", 8)
	body := testutil.Rand("entrybody", 8)
	_, err := fixture.app.EditEntry(ctx, user, blog, entry, title, body)
	assert.NoError(t, err)

	updated, err := fixture.repo.Entry().FindByID(ctx, entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, entry.BlogID, updated.BlogID)
	assert.Equal(t, title, updated.Title)
	assert.Equal(t, body, updated.Body)
	assert.Equal(t, "<p>"+body+"</p>", updated.BodyHTML)
}

func Test_App_EditEntry_NoPermission(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	user2 := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("entrytitle", 8)
	body := testutil.Rand("entrybody", 8)
	_, err := fixture.app.EditEntry(ctx, user2, blog, entry, title, body)
	assert.Equal(t, app.ErrPermissionDenied, err)
}

func Test_App_EditEntry_NotFoundInBlog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	blog2 := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("entrytitle", 8)
	body := testutil.Rand("entrybody", 8)
	_, err := fixture.app.EditEntry(ctx, user, blog2, entry, title, body)
	assert.Equal(t, app.ErrNotFound, err)
}
func Test_App_EditEntry_TooLongTitle(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("", 501)
	body := testutil.Rand("entrybody", 8)
	_, err := fixture.app.EditEntry(ctx, user, blog, entry, title, body)
	assert.Equal(t, app.ErrInvalidArgument, err)
}

func Test_App_UnpublishEntry_Success(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	err := fixture.app.UnpublishEntry(ctx, user, blog, entry)
	assert.NoError(t, err)

	_, err = fixture.repo.Entry().FindByID(ctx, entry.ID)
	assert.Equal(t, app.ErrNotFound, err)
}

func Test_App_UnpublishEntry_NoPermission(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	user2 := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	err := fixture.app.UnpublishEntry(ctx, user2, blog, entry)
	assert.Equal(t, app.ErrPermissionDenied, err)
}

func Test_App_UnpublishEntry_NotFoundInBlog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	blog2 := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)
	ctx := context.Background()

	err := fixture.app.UnpublishEntry(ctx, user, blog2, entry)
	assert.Equal(t, app.ErrNotFound, err)
}
