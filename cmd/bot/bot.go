package bot

import (
	"bm/config"
	"bm/controller"
	"bm/db"
	"bm/di"
	"bm/i18n"
	"bm/repository"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "bot",
		Action: botAction,
	}
}

func botAction(cCtx *cli.Context) error {
	fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.SlogLogger{Logger: slog.Default()}
		}),
		fx.Provide(config.GetConfig),

		fx.Provide(
			di.InitBot,
			fx.Annotate(
				di.InitBotUpdater,
				fx.ParamTags("", "", "", `group:"controllers"`),
			),
		),

		fx.Provide(i18n.NewMessageBundle),
		fx.Provide(db.NewDB),
		fx.Provide(repository.NewUserRepo),

		fx.Provide(di.AsController(controller.NewStartController)),

		fx.Invoke(func(*ext.Updater) {}),
	).Run()

	return nil
}
