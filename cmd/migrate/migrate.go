package migrate

import (
	"bm/config"
	"bm/db"
	"bm/entity"
	"errors"

	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "migrate",
		Action: migrateAction,
	}
}

func migrateAction(cCtx *cli.Context) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	db, err := db.NewDB(cfg)
	if err != nil {
		return err
	}

	return errors.Join(
		db.AutoMigrate(&entity.User{}),
	)
}
