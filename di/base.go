package di

import (
	"bm/config"
	"bm/db"
	"bm/i18n"
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func BaseModule() fx.Option {
	return fx.Module(
		"base",
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.SlogLogger{Logger: slog.Default()}
		}),
		fx.Provide(config.GetConfig),
		fx.Provide(db.NewDB),
		fx.Provide(fx.Annotate(
			db.NewEm,
			fx.As(new(db.Saver)),
		)),
		fx.Provide(i18n.NewMessageBundle),
	)
}
