package app

import (
	"context"
	"testing"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// XXX: 他のテストで作成したブログも検索してしまうためうまくテストできない

// func Test_App_ListBlogs_All(t *testing.T) {
// 	fixture := setup()
// 	user := testutil.CreateUser(fixture.repo)
// 	user2 := testutil.CreateUser(fixture.repo)
// 	blog1 := testutil.CreateBlog(user, fixture.repo)
// 	blog2 := testutil.CreateBlog(user, fixture.repo)
// 	blog3 := testutil.CreateBlog(user2, fixture.repo)
// 	ctx := context.Background()

// 	blogs, hasNextPage, err := fixture.app.ListBlogs(ctx, 1, 5)
// 	assert.NoError(t, err)
// 	assert.Len(t, blogs, 3)
// 	assert.Contains(t, blogs, blog1)
// 	assert.Contains(t, blogs, blog2)
// 	assert.Contains(t, blogs, blog3)
// 	assert.False(t, hasNextPage)
// }

// func Test_App_ListBlogs_Paging(t *testing.T) {
// 	fixture := setup()
// 	user := testutil.CreateUser(fixture.repo)
// 	user2 := testutil.CreateUser(fixture.repo)
// 	testutil.CreateBlog(user, fixture.repo)
// 	testutil.CreateBlog(user, fixture.repo)
// 	testutil.CreateBlog(user2, fixture.repo)
// 	ctx := context.Background()

// 	blogs, hasNextPage, err := fixture.app.ListBlogs(ctx, 1, 2)
// 	assert.NoError(t, err)
// 	assert.Len(t, blogs, 2)
// 	assert.True(t, hasNextPage)

// 	blogs, hasNextPage, err = fixture.app.ListBlogs(ctx, 2, 2)
// 	assert.NoError(t, err)
// 	assert.Len(t, blogs, 1)
// 	assert.False(t, hasNextPage)
// }

func Test_App_ListBlogsByUser_All(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	user2 := testutil.CreateUser(fixture.repo)
	blog1 := testutil.CreateBlog(user, fixture.repo)
	blog2 := testutil.CreateBlog(user, fixture.repo)
	blog3 := testutil.CreateBlog(user, fixture.repo)
	blog4 := testutil.CreateBlog(user2, fixture.repo)
	ctx := context.Background()

	blogs, hasNextPage, err := fixture.app.ListBlogsByUser(ctx, user, 1, 5)
	assert.NoError(t, err)
	assert.Len(t, blogs, 3)
	assert.Contains(t, blogs, blog1)
	assert.Contains(t, blogs, blog2)
	assert.Contains(t, blogs, blog3)
	assert.NotContains(t, blogs, blog4)
	assert.False(t, hasNextPage)
}

func Test_App_ListBlogsByUser_Paging(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	testutil.CreateBlog(user, fixture.repo)
	testutil.CreateBlog(user, fixture.repo)
	testutil.CreateBlog(user, fixture.repo)
	ctx := context.Background()

	blogs, hasNextPage, err := fixture.app.ListBlogsByUser(ctx, user, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, blogs, 2)
	assert.True(t, hasNextPage)

	blogs, hasNextPage, err = fixture.app.ListBlogsByUser(ctx, user, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, blogs, 1)
	assert.False(t, hasNextPage)
}

func Test_App_CreateBlog_Success(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	ctx := context.Background()

	path := testutil.Rand("blogpath", 8)
	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("blogdesc", 8)
	blog, err := fixture.app.CreateBlog(ctx, user, path, title, desc)
	assert.NoError(t, err)

	blog, err = fixture.repo.Blog().FindByID(ctx, blog.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, blog.UserID)
	assert.Equal(t, path, blog.Path)
	assert.Equal(t, title, blog.Title)
	assert.Equal(t, desc, blog.Description)
}

func Test_App_CreateBlog_InvalidPath(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	ctx := context.Background()

	path := testutil.Rand("!/invalid", 8)
	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("blogdesc", 8)
	_, err := fixture.app.CreateBlog(ctx, user, path, title, desc)
	assert.Equal(t, app.ErrInvalidArgument, err)
}

func Test_App_CreateBlog_TooLongTitle(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	ctx := context.Background()

	path := testutil.Rand("blogpath", 8)
	title := testutil.Rand("", 201)
	desc := testutil.Rand("blogdesc", 8)
	_, err := fixture.app.CreateBlog(ctx, user, path, title, desc)
	assert.Equal(t, app.ErrInvalidArgument, err)
}

func Test_App_CreateBlog_TooLongDescription(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	ctx := context.Background()

	path := testutil.Rand("blogpath", 8)
	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("", 501)
	_, err := fixture.app.CreateBlog(ctx, user, path, title, desc)
	assert.Equal(t, app.ErrInvalidArgument, err)
}

func Test_App_FindBlogByPath_Found(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	ctx := context.Background()

	path := testutil.Rand("blogpath", 8)
	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("blogdesc", 8)
	blog, err := fixture.app.CreateBlog(ctx, user, path, title, desc)
	assert.NoError(t, err)

	found, err := fixture.app.FindBlogByPath(ctx, path)
	assert.NoError(t, err)
	assert.Equal(t, blog.ID, found.ID)
}

func Test_App_FindBlogByPath_NotFound(t *testing.T) {
	fixture := setup()
	ctx := context.Background()

	_, err := fixture.app.FindBlogByPath(ctx, "unknown-path")
	assert.Equal(t, app.ErrNotFound, err)
}

func Test_App_EditBlog_Success(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("blogdesc", 8)
	_, err := fixture.app.EditBlog(ctx, user, blog, title, desc)
	assert.NoError(t, err)

	updated, err := fixture.repo.Blog().FindByID(ctx, blog.ID)
	assert.NoError(t, err)
	assert.Equal(t, blog.UserID, updated.UserID)
	assert.Equal(t, blog.Path, updated.Path)
	assert.Equal(t, title, updated.Title)
	assert.Equal(t, desc, updated.Description)
}

func Test_App_EditBlog_NoPermission(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	user2 := testutil.CreateUser(fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("blogdesc", 8)
	_, err := fixture.app.EditBlog(ctx, user2, blog, title, desc)
	assert.Equal(t, app.ErrPermissionDenied, err)
}

func Test_App_EditBlog_TooLongTitle(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("", 201)
	desc := testutil.Rand("blogdesc", 8)
	_, err := fixture.app.EditBlog(ctx, user, blog, title, desc)
	assert.Equal(t, app.ErrInvalidArgument, err)
}

func Test_App_EditBlog_TooLongDescription(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	ctx := context.Background()

	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("", 501)
	_, err := fixture.app.EditBlog(ctx, user, blog, title, desc)
	assert.Equal(t, app.ErrInvalidArgument, err)
}

func Test_App_DeleteBlog_Success(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	ctx := context.Background()

	err := fixture.app.DeleteBlog(ctx, user, blog)
	assert.NoError(t, err)

	_, err = fixture.repo.Blog().FindByID(ctx, blog.ID)
	assert.Equal(t, domain.ErrNotFound, err)
}

func Test_App_DeleteBlog_NoPermission(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	user2 := testutil.CreateUser(fixture.repo)
	ctx := context.Background()

	err := fixture.app.DeleteBlog(ctx, user2, blog)
	assert.Equal(t, app.ErrPermissionDenied, err)
}
