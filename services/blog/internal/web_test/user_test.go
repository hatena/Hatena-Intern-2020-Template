package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_WillSignup_SignedOut(t *testing.T) {
	fixture := setup()

	req := httptest.NewRequest(http.MethodGet, "/signup", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, nil)
	c.SetPath("/signup")
	if assert.NoError(t, fixture.server.WillSignupHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// フォームが表示される
		assert.Contains(t, rec.Body.String(), `<form method="POST" action="/signup">`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="name">`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="password">`)
	}
}

func Test_WillSignup_SignedIn(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/signup", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/signup")
	if assert.NoError(t, fixture.server.WillSignupHandler()(c)) {
		// サインイン済みの場合はリダイレクトされる
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/", rec.Header().Get("Location"))
	}
}

func Test_Signup_Success(t *testing.T) {
	fixture := setup()

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	f := make(url.Values)
	f.Set("name", name)
	f.Set("password", password)
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, nil)
	c.SetPath("/signup")
	if assert.NoError(t, fixture.server.SignupHandler()(c)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/", rec.Header().Get("Location"))
		// サインアップできている
		var key string
		for _, c := range rec.Result().Cookies() {
			if c.Name == "INTERN_2020_BLOG_SESSION" {
				key = c.Value
				break
			}
		}
		assert.NotEmpty(t, key)
		user, _, err := fixture.app.FindUserBySessionKey(context.Background(), key)
		assert.NoError(t, err)
		assert.Equal(t, name, user.Name)
	}
}

func Test_WillSignin_SignedOut(t *testing.T) {
	fixture := setup()

	req := httptest.NewRequest(http.MethodGet, "/signin", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, nil)
	c.SetPath("/signin")
	if assert.NoError(t, fixture.server.WillSigninHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// フォームが表示される
		assert.Contains(t, rec.Body.String(), `<form method="POST" action="/signin">`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="name">`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="password">`)
	}
}

func Test_WillSignin_SignedIn(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/signin", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/signin")
	if assert.NoError(t, fixture.server.WillSigninHandler()(c)) {
		// サインイン済みの場合はリダイレクトされる
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/", rec.Header().Get("Location"))
	}
}

func Test_Signin_Success(t *testing.T) {
	fixture := setup()

	name := testutil.Rand("user", 8)
	password := testutil.Rand("", 16)
	f := make(url.Values)
	f.Set("name", name)
	f.Set("password", password)
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, nil)
	c.SetPath("/signin")
	if assert.NoError(t, fixture.server.SigninHandler()(c)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/", rec.Header().Get("Location"))
		// サインインできている
		var key string
		for _, c := range rec.Result().Cookies() {
			if c.Name == "INTERN_2020_BLOG_SESSION" {
				key = c.Value
				break
			}
		}
		assert.NotEmpty(t, key)
		user, _, err := fixture.app.FindUserBySessionKey(context.Background(), key)
		assert.NoError(t, err)
		assert.Equal(t, name, user.Name)
	}
}

func Test_WillSignout(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/signout", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/signout")
	if assert.NoError(t, fixture.server.WillSignoutHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// フォームが表示される
		assert.Contains(t, rec.Body.String(), `<form method="POST" action="/signout">`)
	}
}

func Test_Signout_Success(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)

	req := httptest.NewRequest(http.MethodPost, "/signout", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/signout")
	if assert.NoError(t, fixture.server.SignoutHandler()(c)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/", rec.Header().Get("Location"))
		// サインアウトできている
		var key string
		for _, c := range rec.Result().Cookies() {
			if c.Name == "INTERN_2020_BLOG_SESSION" {
				key = c.Value
				break
			}
		}
		assert.Empty(t, key)
	}
}
