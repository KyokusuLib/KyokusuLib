package main

import (
	"github.com/lanxre/kyokusulib/internal/app"
	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/handlers"
	"github.com/lanxre/kyokusulib/internal/lib/logger"
	"github.com/lanxre/kyokusulib/internal/repository"
	"github.com/lanxre/kyokusulib/internal/routes"
	service "github.com/lanxre/kyokusulib/internal/services"
	rhService "github.com/lanxre/kyokusulib/internal/parse/service/ranobehub"
	tgService "github.com/lanxre/kyokusulib/internal/services/tg"
	"github.com/lanxre/kyokusulib/internal/sse"
	"github.com/lanxre/kyokusulib/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(
		fx.Provide(
			config.Load,
			logger.New,			
			storage.NewPostgresDB,
			storage.NewRedisClient,
			app.NewValidator,
			

			sse.NewNotificationHub,
			sse.NewTgPostsHub,

			repository.NewNotificationRepository,
			repository.NewTgPostRepository,
			repository.NewUserRepository,
			repository.NewUserProfileSettingRepository,
			repository.NewUserSocialsRepository,
			repository.NewUserActivityRepository,
			repository.NewUserProfileRepository,
			repository.NewAuthorRepository,
			repository.NewNovelaRatingRepository,
			repository.NewNovelaBookmarkRepository,
			repository.NewNovelaLikeRepository,
			repository.NewCommentsRepository,
			repository.NewNovelaRepository,
			repository.NewTeamRepository,
			repository.NewNovelaStatisticsRepository,
			repository.NewUserTagRepository,
			repository.NewCatalogRepository,
			repository.NewUserExperianceRepository,

			service.NewNotificationService,
			tgService.NewTgPostService,
			service.NewAuthService,
			service.NewUserService,
			service.NewUserActivityService,
			service.NewProfileSettingService,
			service.NewSocialsService,
			service.NewAuthorService,
			service.NewNovelaService,
			service.NewCommentService,
			service.NewTeamService,
			service.NewModerationService,
			service.NewCatalogService,
			rhService.NewRanobeHubParseService,
			tgService.NewTgBotAPI,
			service.NewNovelaStatisticsService,
			service.NewUserExperianceService,
			service.NewUserTagService,
			app.NewEmailService,

			handlers.NewNotificationHandler,
			handlers.NewTgPostHandler,
			handlers.NewAuthHandler,
			handlers.NewHealthHandler,
			handlers.NewUserHandler,
			handlers.NewEmailHandler,
			handlers.NewProfileSettingHandler,
			handlers.NewSocialNetworkHandler,
			handlers.NewAuthorHandler,
			handlers.NewUserActivityHandler,
			handlers.NewNovelaHandler,
			handlers.NewCommentHandler,
			handlers.NewTeamHandler,
			handlers.NewModerationHandler,
			handlers.NewParseHandler,
			handlers.NewNovelaStatisticsHandler,
			handlers.NewCatalogHandler,
			handlers.NewUserExperianceHandler,
			handlers.NewUserTagHandler,

			app.AsRoute(func(h *handlers.HealthHandler) *routes.HealthRoutes { return &routes.HealthRoutes{Handler: h} }),
			app.AsRoute(func(h *handlers.AuthHandler) *routes.AuthRoutes { return &routes.AuthRoutes{Handler: h} }),
			app.AsRoute(func(h *handlers.EmailHandler) *routes.EmailRoutes { return &routes.EmailRoutes{Handler: h} }),
			app.AsRoute(func(h *handlers.ProfileSettingHandler) *routes.ProfileSettingRoutes {
				return &routes.ProfileSettingRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.SocialNetworkHandler) *routes.SocialNetworkRoutes {
				return &routes.SocialNetworkRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.UserHandler) *routes.UserRoutes { return &routes.UserRoutes{Handler: h} }),
			app.AsRoute(func(h *handlers.AuthorHandler) *routes.AuthorRoutes { return &routes.AuthorRoutes{Handler: h} }),
			app.AsRoute(func(h *handlers.UserActivityHandler) *routes.UserActivityRoutes {
				return &routes.UserActivityRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.NovelaHandler) *routes.NovelaRoutes { return &routes.NovelaRoutes{Handler: h} }),
			app.AsRoute(func(h *handlers.TeamHandler) *routes.TeamRoutes { return &routes.TeamRoutes{Handler: h} }),
			app.AsRoute(func(h *handlers.CommentHandler) *routes.CommentRoutes { return &routes.CommentRoutes{Handler: h} }),
			app.AsRoute(func(h *handlers.NotificationHandler) *routes.NotificationRoutes {
				return &routes.NotificationRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.TgPostHandler) *routes.TgPostRoutes {
				return &routes.TgPostRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.ModerationHandler) *routes.ModerationRoutes {
				return &routes.ModerationRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.ParseHandler) *routes.ParseRoutes {
				return &routes.ParseRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.NovelaStatisticsHandler) *routes.NovelaStatisticsRoutes {
				return &routes.NovelaStatisticsRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.CatalogHandler) *routes.CatalogRoutes {
				return &routes.CatalogRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.UserExperianceHandler) *routes.UserExperianceRoutes {
				return &routes.UserExperianceRoutes{Handler: h}
			}),
			app.AsRoute(func(h *handlers.UserTagHandler) *routes.UserTagRoutes {
				return &routes.UserTagRoutes{Handler: h}
			}),
			
			fx.Annotate(
				app.NewMuxRouter,
				fx.ParamTags("", `group:"routes"`),
			),
		),

		fx.Invoke(
			app.RegisterStaticFiles,
			app.StartBackgroundWorkers,
			app.RegisterNotificationHubCleanup,
			app.RegisterTgPostsHubCleanup,
			app.RegisterRedisCleanup,
			app.StartHTTPServer,
		),
		fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
	).Run()
}