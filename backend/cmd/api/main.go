package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"paypath/internal/server"
	"paypath/pkg/logger"
	"paypath/pkg/setting"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := setting.Load()
	logger.Init(cfg.Environment)

	srv := server.New(cfg)

	if err := srv.Start(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info().Msg("server stopped")
		} else {
			logger.Log.Fatal().Err(err).Msg("server failed")
		}
	}
}
