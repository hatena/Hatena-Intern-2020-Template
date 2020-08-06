package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Blog(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry1 := testutil.CreateEntry(blog, fixture.repo)
	entry2 := testutil.CreateEntry(blog, fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/blogs/"+blog.Path, nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/blogs/" + blog.Path)
	c.SetParamNames("path")
	c.SetParamValues(blog.Path)
	if assert.NoError(t, fixture.server.BlogHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// ブログの内容が表示される
		assert.Contains(t, rec.Body.String(), blog.Title)
		assert.Contains(t, rec.Body.String(), blog.Description)
		assert.Contains(t, rec.Body.String(), entry1.Title)
		assert.Contains(t, rec.Body.String(), entry2.Title)
	}
}

func Test_Entry(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)
	blog := testutil.CreateBlog(user, fixture.repo)
	entry := testutil.CreateEntry(blog, fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/blogs/"+blog.Path+"/entries/"+entry.ID.String(), nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/blogs/" + blog.Path + "/entries/" + entry.ID.String())
	c.SetParamNames("path", "id")
	c.SetParamValues(blog.Path, entry.ID.String())
	if assert.NoError(t, fixture.server.BlogHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// エントリの内容が表示される
		assert.Contains(t, rec.Body.String(), blog.Title)
		assert.Contains(t, rec.Body.String(), blog.Description)
		assert.Contains(t, rec.Body.String(), entry.Title)
		assert.Contains(t, rec.Body.String(), entry.BodyHTML)
	}
}
