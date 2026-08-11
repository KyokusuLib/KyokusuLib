package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/config"
	service "github.com/lanxre/kyokusulib/internal/services"
	"github.com/lanxre/kyokusulib/internal/sse"
	"github.com/lanxre/kyokusulib/internal/utils/static"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"go.uber.org/fx"
)

func RegisterRedisCleanup(lc fx.Lifecycle, client *redis.Client) {
	if client == nil {
		return
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing redis client...")
			return client.Close()
		},
	})
}

func RegisterStaticFiles(r *mux.Router) {
	static.CreateStaticDirs(r)
}

func StartBackgroundWorkers(lc fx.Lifecycle, authService *service.AuthService, hub *sse.NotificationHub, tgPostsHub *sse.TgPostsHub) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go authService.StartCleanupWorker(ctx)
			hub.Start(ctx)
			tgPostsHub.Start(ctx)
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})
}

func RegisterTgPostsHubCleanup(lc fx.Lifecycle, hub *sse.TgPostsHub) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing tg posts hub...")
			hub.Close()
			return nil
		},
	})
}

func RegisterNotificationHubCleanup(lc fx.Lifecycle, hub *sse.NotificationHub) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing notification hub...")
			hub.Close()
			return nil
		},
	})
}

func StartHTTPServer(lc fx.Lifecycle, r *mux.Router, cfg *config.Config) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "Cookie"},
		AllowCredentials: true,
		Debug:            false,
	})

	handler := c.Handler(r)

	srv := &http.Server{
		Handler:           handler,
		Addr:              cfg.Address,
		WriteTimeout:      cfg.WriteTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Server starting on %s ...", cfg.Address)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return srv.Shutdown(ctx)
		},
	})
}