package closer

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GracefulCloser struct {
	closingFunc []func() error
	timeout     time.Duration
}

func NewGracefulCloser() *GracefulCloser {
	return &GracefulCloser{
		closingFunc: make([]func() error, 0),
		timeout:     30 * time.Second,
	}
}

func (g *GracefulCloser) Add(closingFunc func() error) {
	g.closingFunc = append(g.closingFunc, closingFunc)
}

func (g *GracefulCloser) SetTimeout(timeout time.Duration) {
	g.timeout = timeout
}

func (g *GracefulCloser) Run(ctx context.Context, log *slog.Logger) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	log.Info("GracefulCloser started", slog.Int("functions_count", len(g.closingFunc)))

	select {
	case sig := <-signals:
		log.Info("Received shutdown signal", slog.String("signal", sig.String()))
	case <-ctx.Done():
		log.Info("Context cancelled")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	for i := len(g.closingFunc) - 1; i >= 0; i-- {
		select {
		case <-shutdownCtx.Done():
			log.Warn("Shutdown timeout reached")
			return
		default:
		}

		if err := g.closingFunc[i](); err != nil {
			log.Error("Error during closing function", slog.String("error", err.Error()))
		}
	}

	log.Info("Shutdown completed")
}
