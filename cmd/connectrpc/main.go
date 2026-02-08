package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/caarlos0/env/v11"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"

	"github.com/pcriv/mancala/internal/mancala"
	redisstore "github.com/pcriv/mancala/internal/store/redis"
	"github.com/pcriv/mancala/proto/protoconnect"
)

type envConfig struct {
	Env      string `env:"ENV"                envDefault:"local"`
	LogLevel string `env:"LOG_LEVEL"          envDefault:"debug"`
	RedisURL string `env:"REDIS_URL,required"`
	Address  string `env:"ADDRESS"            envDefault:":50051"`
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

	redisClient, err := newRedisClient(ctx, cfg.RedisURL)
	if err != nil {
		return logger, err
	}
	gameStore := redisstore.NewGameStore(redisClient)

	path, handler := protoconnect.NewServiceHandler(
		handler{
			service: mancala.NewService(gameStore),
		},
		connect.WithInterceptors(
			validate.NewInterceptor(),
		),
	)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	// Use h2c so we can serve HTTP/2 without TLS.
	svr := http.Server{
		Addr:              cfg.Address,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 0,
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Info("starting connect server on address", slog.String("address", cfg.Address))

		if svrErr := svr.ListenAndServe(); svrErr != nil {
			return fmt.Errorf("failed to listen and serve connect service: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		if svrErr := svr.Close(); svrErr != nil {
			return fmt.Errorf("failed to close server: %w", err)
		}

		return nil
	})

	return logger, g.Wait()
}

func newRedisClient(ctx context.Context, url string) (*redis.Client, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
