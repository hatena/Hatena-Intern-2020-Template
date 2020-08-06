package web

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/hatena/Hatena-Intern-2020/services/blog/repository"
	"github.com/hatena/Hatena-Intern-2020/services/blog/web"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

var testDB *sqlx.DB

func run(m *testing.M) int {
	rand.Seed(time.Now().UnixNano())

	db, err := testutil.NewDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	testDB = db

	return m.Run()
}

type testFixture struct {
	server         *web.Server
	app            *app.App
	repo           *repository.Repository
	accountClient  *testutil.TestAccountClient
	rendererClient *testutil.TestRendererClient
}

func setup() *testFixture {
	accountClient := testutil.CreateTestAccountClient()
	rendererClient := testutil.CreateTestRendererClient()
	a := app.NewApp(testDB, accountClient, testutil.AccountECDSAPublicKey, rendererClient)
	repo := repository.NewRepository(testDB)
	server, _ := web.NewServer(a)
	return &testFixture{server, a, repo, accountClient, rendererClient}
}

func (f *testFixture) newCustomContext(req *http.Request, rec *httptest.ResponseRecorder, user *domain.User) *web.CustomContext {
	c := f.server.Echo().NewContext(req, rec)
	return &web.CustomContext{Context: c, User: user}
}

func Test_Index_SignedOut(t *testing.T) {
	fixture := setup()

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, nil)
	c.SetPath("/")
	if assert.NoError(t, fixture.server.IndexHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// サインアップ, サインインのナビゲーションが表示される
		assert.Contains(t, rec.Body.String(), `<a href="/signup">サインアップ</a>`)
		assert.Contains(t, rec.Body.String(), `<a href="/signin">サインイン</a>`)
	}
}

func Test_Index_SignedIn(t *testing.T) {
	fixture := setup()
	user := testutil.CreateUser(fixture.repo)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rec := httptest.NewRecorder()
	c := fixture.newCustomContext(req, rec, user)
	c.SetPath("/")
	if assert.NoError(t, fixture.server.IndexHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// マイページ, サインアウトのナビゲーションが表示される
		assert.Contains(t, rec.Body.String(), `<a href="/my/blogs">マイページ</a>`)
		assert.Contains(t, rec.Body.String(), `<a href="/signout">サインアウト</a>`)
	}
}
