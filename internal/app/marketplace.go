package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hanoys/marketplace-api/auth"
	"github.com/hanoys/marketplace-api/config"
	"github.com/hanoys/marketplace-api/internal/handler"
	"github.com/hanoys/marketplace-api/internal/repository/postgres"
	"github.com/hanoys/marketplace-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

func formConnectionURL(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
}

func createConnectionPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(ctx, formConnectionURL(cfg))
	if err != nil {
		return nil, err
	}

	if err = dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}

func newRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func createTokenProviderConfig(cfg *config.Config) *auth.ProviderConfig {
	return auth.NewProviderConfig(cfg.JWT.AccessTokenExpTime, cfg.JWT.RefreshTokenExpTime, cfg.JWT.SecretKey)
}

func createAdvertisementServiceConfig(cfg *config.Config) *service.AdvertisementServiceConfig {
	return &service.AdvertisementServiceConfig{
		CheckImageIdleTimeout: cfg.App.CheckImageIdleTimeout,
		MinImageWidth:         cfg.App.MinImageWidth,
		MaxImageWidth:         cfg.App.MaxImageWidth,
		MinImageHeight:        cfg.App.MinImageHeight,
		MaxImageHeight:        cfg.App.MaxImageHeight,
	}
}

func printConfig(cfg *config.Config) error {
	jsonConfig, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	log.Printf("Config: \n%s\n", jsonConfig)
	return nil
}

func Run() {
	cfg, err := config.GetConfig(".env.local")
	if err != nil {
		log.Fatalf("load config error: %v\n", err)
	}

	printConfig(cfg)

	connPool, err := createConnectionPool(context.Background(), cfg)
	if err != nil {
		log.Fatalf("unable to establish connection with database: %v\n", err)
	}

	redisClient, err := newRedisClient(context.Background(), cfg)
	if err != nil {
		log.Fatalf("unable to establish connection with redis: %v\n", err)
	}

	serviceRepository := postgres.NewRepositories(connPool)
	tokenProvider := auth.NewProvider(redisClient, createTokenProviderConfig(cfg))
	services := service.NewServices(serviceRepository, tokenProvider, createAdvertisementServiceConfig(cfg))
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
