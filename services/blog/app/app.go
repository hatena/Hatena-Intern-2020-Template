package app

import (
	"context"
	"crypto/ecdsa"

	"github.com/jmoiron/sqlx"

	pb_account "github.com/hatena/Hatena-Intern-2020/services/blog/pb/account"
	pb_renderer "github.com/hatena/Hatena-Intern-2020/services/blog/pb/renderer"
)

// App はアプリケーションを表す
type App struct {
	db                    *sqlx.DB
	accountClient         pb_account.AccountClient
	accountECDSAPublicKey *ecdsa.PublicKey
	rendererClient        pb_renderer.RendererClient
}

// NewApp は App を作成する
func NewApp(
	db *sqlx.DB,
	accountClient pb_account.AccountClient,
	accountECDSAPublicKey *ecdsa.PublicKey,
	rendererClient pb_renderer.RendererClient,
) *App {
	return &App{db, accountClient, accountECDSAPublicKey, rendererClient}
}

// Render は RendererClient を使った domain.BodyRenderer の実装
func (a *App) Render(ctx context.Context, body string) (string, error) {
	reply, err := a.rendererClient.Render(ctx, &pb_renderer.RenderRequest{Src: body})
	if err != nil {
		return "", err
	}
	return reply.Html, nil
}
