package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/bufbuild/protovalidate-go"
	"github.com/caarlos0/env/v11"
	"github.com/pcriv/mancala/internal/mancala"
	"github.com/pcriv/mancala/proto"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

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

	grpcServer := grpc.NewServer()
	proto.RegisterServiceServer(grpcServer, serviceServer{
		service:   mancala.NewService(gameStore),
		validator: validator,
	})

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", cfg.Address)
		if err != nil {
			return fmt.Errorf("failed to listen on address %q: %w", cfg.Address, err)
		}

		slog.Info("starting grpc server on address", slog.String("address", cfg.Address))

		err = grpcServer.Serve(lis)
		if err != nil {
			return fmt.Errorf("failed to serve grpc service: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		grpcServer.GracefulStop()

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
