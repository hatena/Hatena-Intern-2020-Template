package grpc

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/account/app"
	"github.com/hatena/Hatena-Intern-2020/services/account/grpc"
	"github.com/hatena/Hatena-Intern-2020/services/account/internal/testutil"
	"github.com/jmoiron/sqlx"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
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

func readECDSAPrivateKey(data string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(data))
	key, _ := x509.ParseECPrivateKey(block.Bytes)
	return key
}

func readECDSAPublicKey(data string) *ecdsa.PublicKey {
	block, _ := pem.Decode([]byte(data))
	rawKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	return rawKey.(*ecdsa.PublicKey)
}

// テスト用の鍵
var ecdsaPrivateKey = readECDSAPrivateKey(`
-----BEGIN EC PRIVATE KEY-----
MGgCAQEEHG2LBwHQA+1bhps5+CoylSMFZtGSZ0Ldxwrx9EGgBwYFK4EEACGhPAM6
AAROulH/IepJoaxv8zgiQ6p7LPtOJW+BPESrnHw9ziGICcIjwtrUS/xE+qLC83Tk
+JR+P9oKZHpfXw==
-----END EC PRIVATE KEY-----
`)
var ecdsaPublicKey = readECDSAPublicKey(`
-----BEGIN PUBLIC KEY-----
ME4wEAYHKoZIzj0CAQYFK4EEACEDOgAETrpR/yHqSaGsb/M4IkOqeyz7TiVvgTxE
q5x8Pc4hiAnCI8La1Ev8RPqiwvN05PiUfj/aCmR6X18=
-----END PUBLIC KEY-----
`)

func setup() (*grpc.Server, *app.App) {
	app := app.NewApp(testDB)
	server := grpc.NewServer(&grpc.Config{
		App:             app,
		ECDSAPrivateKey: ecdsaPrivateKey,
	})
	return server, app
}

func verifyToken(src string) (jwt.Token, error) {
	token, err := jwt.Parse(strings.NewReader(src), jwt.WithVerify(jwa.ES256, ecdsaPublicKey))
	if err != nil {
		return nil, err
	}
	return token, nil
}
