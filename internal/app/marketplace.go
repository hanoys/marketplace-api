package app

import (
	"context"
	"fmt"
	"github.com/hanoys/marketplace-api/auth"
	"github.com/hanoys/marketplace-api/internal/config"
	"github.com/hanoys/marketplace-api/internal/handler"
	"github.com/hanoys/marketplace-api/internal/repository/postgres"
	"github.com/hanoys/marketplace-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

func createConnectionPool(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	if err = dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}

func newRedisClient(ctx context.Context, host string, port string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func createTokenProviderConfig(cfg *config.Config) *auth.ProviderConfig {
	return auth.NewProviderConfig(cfg.JWT.AccessTokenExpTime, cfg.JWT.RefreshTokenExpTime, cfg.JWT.SecretKey)
}

func Run() {
	cfg, err := config.GetConfig(".env.local")
	if err != nil {
		log.Fatalf("load config error: %v\n", err)
	}

	fmt.Printf("config: %v\n", cfg)

	//TODO:  URL -> URI
	connPool, err := createConnectionPool(context.Background(), cfg.DB.URL)
	if err != nil {
		log.Fatalf("unable to establish connection with database: %v\n", err)
	}

	redisClient, err := newRedisClient(context.Background(), cfg.Redis.Host, cfg.Redis.Port)
	if err != nil {
		log.Fatalf("unable to establish connection with redis: %v\n", err)
	}

	serviceRepository := postgres.NewRepositories(connPool)
	tokenProvider := auth.NewProvider(redisClient, createTokenProviderConfig(cfg))
	services := service.NewServices(serviceRepository, tokenProvider)
	serviceHandler := handler.NewHandler(services)

	server := http.Server{
		Handler:      serviceHandler.Init(),
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	log.Printf("Starting server at: %v\n", server.Addr)
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
