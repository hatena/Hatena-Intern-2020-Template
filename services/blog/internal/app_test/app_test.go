package app

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/blog/app"
	"github.com/hatena/Hatena-Intern-2020/services/blog/internal/testutil"
	"github.com/hatena/Hatena-Intern-2020/services/blog/repository"
	"github.com/jmoiron/sqlx"
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
	return &testFixture{a, repo, accountClient, rendererClient}
}
