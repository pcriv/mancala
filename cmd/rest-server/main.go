package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v11"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pcriv/mancala/cmd/rest-server/internal/web"
	"github.com/pcriv/mancala/internal/mancala"
	"github.com/pcriv/mancala/internal/openapi"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	oapicodegenmiddleware "github.com/oapi-codegen/echo-middleware"
	redisstore "github.com/pcriv/mancala/internal/store/redis"
	slogecho "github.com/samber/slog-echo"
)

type envConfig struct {
	Env      string `env:"ENV" envDefault:"local"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	RedisURL string `env:"REDIS_URL,required"`
	Port     string `env:"PORT" envDefault:"1323"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	err := run(ctx)
	if err != nil {
		slog.Error("failed to start server", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg := envConfig{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	redisClient, err := newRedisClient(ctx, cfg.RedisURL)
	if err != nil {
		return err
	}
	gameStore := redisstore.NewGameStore(redisClient)

	spec, err := openapi.GetSwagger()
	if err != nil {
		return err
	}

	reqLoggerCfg := slogecho.Config{
		WithSpanID:  true,
		WithTraceID: true,
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(slogecho.NewWithConfig(logger, reqLoggerCfg))
	e.Use(middleware.RequestID())

	options := oapicodegenmiddleware.Options{
		SilenceServersWarning: true,
	}

	e.Use(oapicodegenmiddleware.OapiRequestValidatorWithOptions(spec, &options))
	e.Use(middleware.Recover())

	s := web.NewAPI(spec, mancala.NewService(gameStore))
	openapi.RegisterHandlers(e, s)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return e.Start(":" + cfg.Port)
	})

	eg.Go(func() error {
		<-ctx.Done()
		return e.Shutdown(context.Background())
	})

	return eg.Wait()
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
