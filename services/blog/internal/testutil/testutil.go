package testutil

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL ドライバを使うために必要
	"github.com/hatena/Hatena-Intern-2020/services/blog/db"
	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	pb_account "github.com/hatena/Hatena-Intern-2020/services/blog/pb/account"
	pb_renderer "github.com/hatena/Hatena-Intern-2020/services/blog/pb/renderer"
	"github.com/jmoiron/sqlx"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"google.golang.org/grpc"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewDB はテスト用に DB ハンドルを作成する
func NewDB() (*sqlx.DB, error) {
	databaseDSN := os.Getenv("TEST_DATABASE_DSN")
	if databaseDSN == "" {
		return nil, errors.New("TEST_DATABASE_DSN is not set")
	}
	db, err := db.Connect(databaseDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}

var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// Rand はランダムな文字列を返す
func Rand(prefix string, size int) string {
	key := make([]rune, size)
	for i := range key {
		key[i] = letters[rand.Intn(len(letters))]
	}
	return prefix + string(key)
}

// RandUint はランダムな整数を返す
func RandUint() uint64 {
	return rand.Uint64()
}

// TestRenderer はテスト用の domain.BodyRenderer の実装
type TestRenderer struct {
	render func(body string) (string, error)
}

// CreateTestRenderer は TestRenderer を作成する
func CreateTestRenderer(render func(body string) (string, error)) *TestRenderer {
	return &TestRenderer{render}
}

// Render は render 関数を使ってレンダリングを行う
func (r *TestRenderer) Render(ctx context.Context, body string) (string, error) {
	return r.render(body)
}

// CreateUser はユーザーを作成する
func CreateUser(repo domain.Repository) *domain.User {
	accountID := domain.AccountID(RandUint())
	name := Rand("user", 8)
	user, _ := domain.CreateUser(accountID, name)(context.Background(), repo)
	return user
}

// CreateBlog はブログを作成する
func CreateBlog(user *domain.User, repo domain.Repository) *domain.Blog {
	path := Rand("blogpath", 8)
	title := Rand("blogtitle", 8)
	desc := Rand("blogdesc", 8)
	blog, _ := user.CreateBlog(path, title, desc)(context.Background(), repo)
	return blog
}

// CreateEntry はエントリを作成する
func CreateEntry(blog *domain.Blog, repo domain.Repository) *domain.Entry {
	title := Rand("entrytitle", 8)
	body := Rand("entrybody", 8)
	publishedAt := time.Now().Truncate(time.Millisecond).UTC()
	renderer := CreateTestRenderer(func(body string) (string, error) {
		return "<p>" + body + "</p>", nil
	})
	entry, _ := blog.PublishEntry(title, body, publishedAt)(context.Background(), repo, renderer)
	return entry
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

// AccountECDSAPrivateKey はテスト用のアカウントサービスで使われる秘密鍵
var AccountECDSAPrivateKey = readECDSAPrivateKey(`
-----BEGIN EC PRIVATE KEY-----
MGgCAQEEHG2LBwHQA+1bhps5+CoylSMFZtGSZ0Ldxwrx9EGgBwYFK4EEACGhPAM6
AAROulH/IepJoaxv8zgiQ6p7LPtOJW+BPESrnHw9ziGICcIjwtrUS/xE+qLC83Tk
+JR+P9oKZHpfXw==
-----END EC PRIVATE KEY-----
`)

// AccountECDSAPublicKey はテスト用のアカウントサービスで使われる公開鍵
var AccountECDSAPublicKey = readECDSAPublicKey(`
-----BEGIN PUBLIC KEY-----
ME4wEAYHKoZIzj0CAQYFK4EEACEDOgAETrpR/yHqSaGsb/M4IkOqeyz7TiVvgTxE
q5x8Pc4hiAnCI8La1Ev8RPqiwvN05PiUfj/aCmR6X18=
-----END PUBLIC KEY-----
`)

func generateToken(id, name string) ([]byte, error) {
	claims := jwt.New()
	now := time.Now()
	iss := "hatena-intern-2020-account"
	if err := claims.Set(jwt.IssuerKey, iss); err != nil {
		return nil, err
	}
	sub := "user"
	if err := claims.Set(jwt.SubjectKey, sub); err != nil {
		return nil, err
	}
	exp := now.Add(time.Hour)
	if err := claims.Set(jwt.ExpirationKey, exp); err != nil {
		return nil, err
	}
	if err := claims.Set("user_id", id); err != nil {
		return nil, err
	}
	if err := claims.Set("user_name", name); err != nil {
		return nil, err
	}
	iat := now
	if err := claims.Set(jwt.IssuedAtKey, iat); err != nil {
		return nil, err
	}
	token, err := jwt.Sign(claims, jwa.ES256, AccountECDSAPrivateKey)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// TestAccountClient はテスト用のアカウントサービスのクライアント
type TestAccountClient struct {
	Error error
}

// CreateTestAccountClient は TestAccountClient を作成する
func CreateTestAccountClient() *TestAccountClient {
	return &TestAccountClient{}
}

// Signup はサインアップ処理のフェイク実装 (Error が指定されない限り常に成功する)
func (c *TestAccountClient) Signup(ctx context.Context, in *pb_account.SignupRequest, opts ...grpc.CallOption) (*pb_account.SignupReply, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	id := strconv.FormatUint(uint64(RandUint()), 10)
	token, err := generateToken(id, in.Name)
	if err != nil {
		return nil, err
	}
	return &pb_account.SignupReply{Token: string(token)}, nil
}

// Signin はサインアップ処理のフェイク実装 (Error が指定されない限り常に成功する)
func (c *TestAccountClient) Signin(ctx context.Context, in *pb_account.SigninRequest, opts ...grpc.CallOption) (*pb_account.SigninReply, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	id := strconv.FormatUint(uint64(RandUint()), 10)
	token, err := generateToken(id, in.Name)
	if err != nil {
		return nil, err
	}
	return &pb_account.SigninReply{Token: string(token)}, nil
}

// TestRendererClient はテスト用の記法変換サービスのクライアント
type TestRendererClient struct {
	Error      error
	RenderFunc func(src string) string
}

// CreateTestRendererClient は TestRendererClient を作成する
func CreateTestRendererClient() *TestRendererClient {
	return &TestRendererClient{}
}

// Render は記法変換のフェイク実装 (Error が指定されない限り常に成功する. RenderFunc が指定されればそれを使う)
func (c *TestRendererClient) Render(ctx context.Context, in *pb_renderer.RenderRequest, opts ...grpc.CallOption) (*pb_renderer.RenderReply, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	html := in.Src
	if c.RenderFunc != nil {
		html = c.RenderFunc(in.Src)
	}
	return &pb_renderer.RenderReply{Html: html}, nil
}
