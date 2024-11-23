package graph

import (
	"bm/api/resolver"
	"bm/di"
	"net/http"

	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "graph",
		Action: graphAction,
	}
}

func graphAction(cCtx *cli.Context) error {
	fx.New(
		di.BaseModule(),
		di.PrometheusMetricsServer(),
		di.Repositories(),

		fx.Provide(resolver.NewResolver),
		fx.Provide(di.InitGraphQL),

		fx.Invoke(func(*http.Server) {}),
	).Run()

	return nil
}
