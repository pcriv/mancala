package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v11"
	"golang.org/x/sync/errgroup"

	"github.com/pcriv/mancala/proto/protoconnect"
)

type envConfig struct {
	Env           string `env:"ENV"             envDefault:"local"`
	LogLevel      string `env:"LOG_LEVEL"       envDefault:"debug"`
	Address       string `env:"ADDRESS"         envDefault:":8080"`
	ConnectRPCURL string `env:"CONNECT_RPC_URL" envDefault:"http://localhost:50051"`
}

func main() {
	logger, err := run()
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("error running application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("closing server gracefully")
}

func run() (*slog.Logger, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cfg := envConfig{}
	if err := env.Parse(&cfg); err != nil {
		return logger, fmt.Errorf("unable to parse config: %w", err)
	}

	client := protoconnect.NewServiceClient(http.DefaultClient, cfg.ConnectRPCURL)

	h := newHandler(client, logger)

	mux := http.NewServeMux()
	h.registerRoutes(mux)

	svr := http.Server{
		Addr:              cfg.Address,
		Handler:           mux,
		ReadHeaderTimeout: 0,
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Info("starting web server", slog.String("address", cfg.Address))

		if svrErr := svr.ListenAndServe(); svrErr != nil && !errors.Is(svrErr, http.ErrServerClosed) {
			return fmt.Errorf("failed to listen and serve web server: %w", svrErr)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		if svrErr := svr.Shutdown(context.Background()); svrErr != nil {
			return fmt.Errorf("failed to shutdown server: %w", svrErr)
		}

		return nil
	})

	return logger, g.Wait()
}
