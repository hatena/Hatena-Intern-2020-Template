package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/hatena/Hatena-Intern-2020/services/renderer-go/config"
	server "github.com/hatena/Hatena-Intern-2020/services/renderer-go/grpc"
	"github.com/hatena/Hatena-Intern-2020/services/renderer-go/log"
	pb "github.com/hatena/Hatena-Intern-2020/services/renderer-go/pb/renderer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	// 設定をロード
	conf, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %+v", err)
	}

	// ロガーを初期化
	logger, err := log.NewLogger(log.Config{Mode: conf.Mode})
	if err != nil {
		return fmt.Errorf("failed to create logger: %+v", err)
	}
	defer logger.Sync()

	// サーバーを起動
	logger.Info(fmt.Sprintf("starting gRPC server (port = %v)", conf.GRPCPort))
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(conf.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %+v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(
				logger,
				grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
					// ヘルスチェックのログを無効化
					return !strings.HasPrefix(fullMethodName, "/grpc.health.v1.Health/")
				}),
			),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	svr := server.NewServer()
	pb.RegisterRendererServer(s, svr)
	healthpb.RegisterHealthServer(s, svr)
	go stop(s, conf.GracefulStopTimeout, logger)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %+v", err)
	}

	return nil
}

func stop(s *grpc.Server, timeout time.Duration, logger *zap.Logger) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	logger.Info(fmt.Sprintf("gracefully stopping server (sig = %v)", sig))
	t := time.NewTimer(timeout)
	defer t.Stop()
	stopped := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(stopped)
	}()
	select {
	case <-t.C:
		logger.Warn(fmt.Sprintf("stopping server (not stopped in %s)", timeout.String()))
		s.Stop()
	case <-stopped:
	}
}
