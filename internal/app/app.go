package app

import (
	"log/slog"
	"time"

	"github.com/Nikita-Smirnov-idk/go_microservices_template_project/services/sso/v1/internal/transport/grpc"
)

type App struct {
	GRPCSrv *grpc.Server
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	grpcApp := grpc.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
