package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_MyBlogs(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog1 := testutil.CreateBlog(user, fixture.repo)
	blog2 := testutil.CreateBlog(user, fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/my/blogs", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs")
	if assert.NoError(t, fixture.server.MyBlogsHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// ブログの一覧が表示される
		assert.Contains(t, rec.Body.String(), blog1.Title)
		assert.Contains(t, rec.Body.String(), blog2.Title)
	}
}

func Test_WillCreateBlog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/my/blogs/-/create", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/-/create")
	if assert.NoError(t, fixture.server.WillCreateBlogHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// フォームが表示される
		assert.Contains(t, rec.Body.String(), `<form method="POST" action="/my/blogs">`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="path">`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="title">`)
		assert.Contains(t, rec.Body.String(), `<textarea name="description"></textarea>`)
	}
}

func Test_CreateBlog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)

	path := testutil.Rand("blogpath", 8)
	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("blogdesc", 8)
	f := make(url.Values)
	f.Set("path", path)
	f.Set("title", title)
	f.Set("description", desc)
	req := httptest.NewRequest(http.MethodPost, "/my/blogs", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs")
	if assert.NoError(t, fixture.server.CreateBlogHandler()(c)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/my/blogs/"+path, rec.Header().Get("Location"))
		// ブログが作成されている
		blog, err := fixture.app.FindBlogByPath(context.Background(), path)
		assert.NoError(t, err)
		assert.Equal(t, title, blog.Title)
		assert.Equal(t, desc, blog.Description)
	}
}

func Test_MyBlog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry1 := testutil.CreateEntry(blog, fixture.repo)
	entry2 := testutil.CreateEntry(blog, fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/my/blogs/"+blog.Path, nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/" + blog.Path)
	c.SetParamNames("path")
	c.SetParamValues(blog.Path)
	if assert.NoError(t, fixture.server.MyBlogHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// ブログの内容が表示される
		assert.Contains(t, rec.Body.String(), blog.Title)
		assert.Contains(t, rec.Body.String(), blog.Description)
		assert.Contains(t, rec.Body.String(), entry1.Title)
		assert.Contains(t, rec.Body.String(), entry2.Title)
	}
}

func Test_WillEditBlog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/my/blogs/"+blog.Path+"/edit", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/" + blog.Path + "/edit")
	c.SetParamNames("path")
	c.SetParamValues(blog.Path)
	if assert.NoError(t, fixture.server.WillEditBlogHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// フォームが表示される
		assert.Contains(t, rec.Body.String(), `<form method="POST" action="/my/blogs/`+blog.Path+`/edit">`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="title" value="`+blog.Title+`">`)
		assert.Contains(t, rec.Body.String(), `<textarea name="description">`+blog.Description+`</textarea>`)
	}
}

func Test_EditBlog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)

	title := testutil.Rand("blogtitle", 8)
	desc := testutil.Rand("blogdesc", 8)
	f := make(url.Values)
	f.Set("title", title)
	f.Set("description", desc)
	req := httptest.NewRequest(http.MethodPost, "/my/blogs/"+blog.Path+"/edit", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/" + blog.Path + "/edit")
	c.SetParamNames("path")
	c.SetParamValues(blog.Path)
	if assert.NoError(t, fixture.server.EditBlogHandler()(c)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/my/blogs/"+blog.Path, rec.Header().Get("Location"))
		// ブログが更新されている
		blog, err := fixture.app.FindBlogByPath(context.Background(), blog.Path)
		assert.NoError(t, err)
		assert.Equal(t, title, blog.Title)
		assert.Equal(t, desc, blog.Description)
	}
}

func Test_DeleteBlog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)

	req := httptest.NewRequest(http.MethodPost, "/my/blogs/"+blog.Path+"/delete", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/" + blog.Path + "/delete")
	c.SetParamNames("path")
	c.SetParamValues(blog.Path)
	if assert.NoError(t, fixture.server.DeleteBlogHandler()(c)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/my/blogs", rec.Header().Get("Location"))
		// ブログが削除されている
		_, err := fixture.app.FindBlogByPath(context.Background(), blog.Path)
		assert.Equal(t, app.ErrNotFound, err)
	}
}

func Test_WillPublishEntry(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/my/blogs/"+blog.Path+"/entries/-/publish", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/" + blog.Path + "/entries/-/publish")
	c.SetParamNames("path")
	c.SetParamValues(blog.Path)
	if assert.NoError(t, fixture.server.WillPublishEntryHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// フォームが表示される
		assert.Contains(t, rec.Body.String(), `<form method="POST" action="/my/blogs/`+blog.Path+`/entries/-/publish">`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="title">`)
		assert.Contains(t, rec.Body.String(), `<textarea name="body"></textarea>`)
	}
}

func Test_PublishEntry(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)

	title := testutil.Rand("entrytitle", 8)
	body := testutil.Rand("entrybody", 8)
	f := make(url.Values)
	f.Set("title", title)
	f.Set("body", body)
	req := httptest.NewRequest(http.MethodPost, "/my/blogs/"+blog.Path+"/entries/-/publish", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/" + blog.Path + "/entries/-/publish")
	c.SetParamNames("path")
	c.SetParamValues(blog.Path)
	if assert.NoError(t, fixture.server.PublishEntryHandler()(c)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/my/blogs/"+blog.Path, rec.Header().Get("Location"))
		// エントリが投稿されている
		entries, _, err := fixture.app.ListEntriesByBlog(context.Background(), blog, 1, 2)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(entries))
		assert.Equal(t, title, entries[0].Title)
		assert.Equal(t, body, entries[0].Body)
	}
}

func Test_MyEntry(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/my/blogs/"+blog.Path+"/entries/"+entry.ID.String(), nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/" + blog.Path + "/entries/" + entry.ID.String())
	c.SetParamNames("path", "id")
	c.SetParamValues(blog.Path, entry.ID.String())
	if assert.NoError(t, fixture.server.MyEntryHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// フォームが表示される
		assert.Contains(t, rec.Body.String(), `<form method="POST" action="/my/blogs/`+blog.Path+`/entries/`+entry.ID.String()+`/edit">`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="title" value="`+entry.Title+`">`)
		assert.Contains(t, rec.Body.String(), `<textarea name="body">`+entry.Body+`</textarea>`)
		assert.Contains(t, rec.Body.String(), `<form method="POST" action="/my/blogs/`+blog.Path+`/entries/`+entry.ID.String()+`/unpublish">`)
	}
}

func Test_EditEntry(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)

	title := testutil.Rand("entrytitle", 8)
	body := testutil.Rand("entrybody", 8)
	f := make(url.Values)
	f.Set("title", title)
	f.Set("body", body)
	req := httptest.NewRequest(http.MethodPost, "/my/blogs/"+blog.Path+"/entries/"+entry.ID.String()+"/edit", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/" + blog.Path + "/entries/" + entry.ID.String() + "/edit")
	c.SetParamNames("path", "id")
	c.SetParamValues(blog.Path, entry.ID.String())
	if assert.NoError(t, fixture.server.EditEntryHandler()(c)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/my/blogs/"+blog.Path, rec.Header().Get("Location"))
		// エントリが更新されている
		updated, err := fixture.app.FindEntryByID(context.Background(), blog, entry.ID)
		assert.NoError(t, err)
		assert.Equal(t, title, updated.Title)
		assert.Equal(t, body, updated.Body)
	}
}

func Test_UnpublishEntry(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)

	req := httptest.NewRequest(http.MethodPost, "/my/blogs/"+blog.Path+"/entries/"+entry.ID.String()+"/unpublish", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/my/blogs/" + blog.Path + "/entries/" + entry.ID.String() + "/unpublish")
	c.SetParamNames("path", "id")
	c.SetParamValues(blog.Path, entry.ID.String())
	if assert.NoError(t, fixture.server.UnpublishEntryHandler()(c)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/my/blogs/"+blog.Path, rec.Header().Get("Location"))
		// エントリが削除されている
		_, err := fixture.app.FindEntryByID(context.Background(), blog, entry.ID)
		assert.Equal(t, app.ErrNotFound, err)
	}
}
