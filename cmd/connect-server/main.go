package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/bufbuild/protovalidate-go"
	"github.com/caarlos0/env/v11"
	"github.com/pcriv/mancala/internal/mancala"
	"github.com/pcriv/mancala/proto/protoconnect"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"

	redisstore "github.com/pcriv/mancala/internal/store/redis"
)

type envConfig struct {
	Env      string `env:"ENV" envDefault:"local"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	RedisURL string `env:"REDIS_URL,required"`
	Address  string `env:"ADDRESS" envDefault:":50051"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	err := run(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("error running application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("closing server gracefully")
}

func run(ctx context.Context) error {
	cfg := envConfig{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	redisClient, err := newRedisClient(ctx, cfg.RedisURL)
	if err != nil {
		return err
	}
	gameStore := redisstore.NewGameStore(redisClient)

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}

	path, handler := protoconnect.NewServiceHandler(
		handler{
			service:   mancala.NewService(gameStore),
			validator: validator,
		},
	)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	// Use h2c so we can serve HTTP/2 without TLS.
	srv := http.Server{
		Addr:    cfg.Address,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		slog.Info("starting connect server on address", slog.String("address", cfg.Address))

		if err := srv.ListenAndServe(); err != nil {
			return fmt.Errorf("failed to listen and serve connect service: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		if err := srv.Close(); err != nil {
			return fmt.Errorf("failed to close server: %w", err)
		}

		return nil
	})

	return g.Wait()
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
