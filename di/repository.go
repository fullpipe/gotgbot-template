package di

import (
	"bm/repository"

	"go.uber.org/fx"
)

func Repositories() fx.Option {
	return fx.Module(
		"Repositories",
		fx.Provide(
			repository.NewUserRepo,
		),
	)
}
