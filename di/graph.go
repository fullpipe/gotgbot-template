package di

import (
	"bm/api/auth"
	"bm/api/directive"
	"bm/api/gen"
	"bm/api/resolver"
	"bm/config"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/99designs/gqlgen-contrib/prometheus"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/rs/cors"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

const mb int64 = 1 << 20

func InitGraphQL(lc fx.Lifecycle, cfg config.Config, db *gorm.DB, resolver *resolver.Resolver) *http.Server {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByRealIP(100, time.Minute))
	r.Use(cors.New(cors.Options{
		// AllowedOrigins: []string{
		// 	"http://localhost:8080",
		// 	"http://localhost:8100",
		// 	"http://localhost:4200",
		// },
		AllowCredentials: false,
		AllowedHeaders:   []string{"*"},
		Debug:            cfg.Debug,
	}).Handler)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(auth.Middleware(cfg.Token, db))
	r.Use(middleware.Recoverer)

	srv := handler.NewDefaultServer(gen.NewExecutableSchema(gen.Config{
		Resolvers: resolver,
		Directives: gen.DirectiveRoot{
			HasRole:  directive.HasRole,
			Validate: directive.Validate,
		},
	}))

	srv.AddTransport(transport.MultipartForm{
		MaxMemory:     32 * mb,
		MaxUploadSize: 50 * mb,
	})

	prometheus.Register()
	srv.Use(prometheus.Tracer{})

	if cfg.Debug {
		r.Handle("/playground", playground.Handler("Bot", "/query"))
		srv.Use(extension.Introspection{})
	}

	r.Handle("/query", srv)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.GraphServer.Host, cfg.GraphServer.Port),
		Handler: r,
	}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		ln, err := net.Listen("tcp", s.Addr)
		if err != nil {
			return err
		}

		slog.Info("starting GraphQL server", slog.Any("addr", s.Addr))

		go s.Serve(ln)

		return nil
	}))

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		return s.Shutdown(ctx)
	}))

	return s
}
