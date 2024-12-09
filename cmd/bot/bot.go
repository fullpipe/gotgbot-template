package bot

import (
	"bm/controller"
	"bm/di"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "bot",
		Action: botAction,
	}
}

func botAction(cCtx *cli.Context) error {
	fx.New(
		di.BaseModule(),
		di.PrometheusMetricsServer(),
		di.Repositories(),

		fx.Provide(
			di.InitBot,
			fx.Annotate(
				di.InitBotUpdater,
				fx.ParamTags("", "", "", `group:"controllers"`),
			),
		),

		fx.Provide(di.AsController(controller.NewStartController)),

		fx.Invoke(func(*ext.Updater) {}),
	).Run()

	return nil
}
