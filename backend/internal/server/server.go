package server

import (
	"context"
	"errors"
	"net/http"

	"paypath/internal/api/router"
	"paypath/internal/config"
	"paypath/internal/storage"
	"paypath/pkg/logger"
	"paypath/pkg/setting"
)

type Server struct {
	httpServer *http.Server
	db         *storage.DB
	cfg        setting.Config
}

func New(cfg setting.Config) *Server {
	app := config.Setup(cfg)
	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.HTTPAddr,
			Handler:           router.New(app.Deps),
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
		db:  app.DB,
		cfg: cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	defer s.db.Close()

	errCh := make(chan error, 1)
	go func() {
		logger.Log.Info().Str("addr", s.httpServer.Addr).Msg("PayPath API listening")
		errCh <- s.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		return s.httpServer.Shutdown(shutCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
