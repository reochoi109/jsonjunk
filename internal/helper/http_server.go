package helper

import (
	"context"
	logger "jsonjunk/pkg/logging"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func GracefulShutdown(ctx context.Context, srv *http.Server) {
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Log.Error("server shutdown failed", zap.Error(err))
		} else {
			logger.Log.Info("server shutdown completed")
		}
	}()
}
