package grpc

import (
	"crypto/ecdsa"

	"github.com/hatena/Hatena-Intern-2020/services/account/app"
	pb "github.com/hatena/Hatena-Intern-2020/services/account/pb/account"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Config はサーバーの設定
type Config struct {
	App             *app.App
	ECDSAPrivateKey *ecdsa.PrivateKey
}

// Server は pb.AccountServer に対する実装
type Server struct {
	pb.UnimplementedAccountServer
	healthpb.UnimplementedHealthServer

	app             *app.App
	ecdsaPrivateKey *ecdsa.PrivateKey
}

// NewServer は gRPC サーバーを作成する
func NewServer(conf *Config) *Server {
	return &Server{
		app:             conf.App,
		ecdsaPrivateKey: conf.ECDSAPrivateKey,
	}
}
