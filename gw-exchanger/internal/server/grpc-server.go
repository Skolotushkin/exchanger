package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Skolotushkin/proto-exchange/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gw-exchanger/internal/config"
	"gw-exchanger/internal/service"
)

func RunGRPCServer(ctx context.Context, exchanger *service.ExchangerService, cfg *config.Config, logger *zap.Logger) error {
	addr := fmt.Sprintf("%s:%s", cfg.GRPCHost, cfg.GRPCPort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcServer, exchanger)

	logger.Info("🚀 gRPC server started", zap.String("addr", addr))

	go func() {
		<-ctx.Done()
		logger.Warn("⏳ shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(lis)
}
