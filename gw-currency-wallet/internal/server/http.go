package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	config "gw-currency-wallet/internal/config"
)

func SetupAndRunServer(cfg *config.Config, handler http.Handler, logger *zap.Logger) {
	srv := &http.Server{
		Addr:           cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:        handler,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
		IdleTimeout:    cfg.Server.IdleTimeout,
	}

	// запуск в отдельной горутине
	go func() {
		logger.Info("starting HTTP server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("could not start server", zap.Error(err))
		}
	}()

	// ожидание сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server gracefully stopped")
}
