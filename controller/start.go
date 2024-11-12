package controller

import (
	"bm/di"
	"bm/repository"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/fullpipe/icu-mf/mf"
)

type StartController struct {
	baseController

	userRepo *repository.UserRepo
}

func NewStartController(
	userRepo *repository.UserRepo,
	mfBundle mf.Bundle,
) *StartController {
	return &StartController{
		userRepo: userRepo,

		baseController: baseController{
			mfBundle: mfBundle,
			userRepo: userRepo,
		},
	}
}

func (c *StartController) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	user := c.User(ctx)

	_, err := b.SendMessage(
		ctx.EffectiveSender.ChatId,
		c.Trans(ctx, "start.intro", mf.Arg("username", user.Username)),
		&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML},
	)

	return err
}

func (c *StartController) Register(bot *gotgbot.Bot, dispatcher *ext.Dispatcher) error {
	dispatcher.AddHandler(handlers.NewCommand("start", c.Start))

	return nil
}

var _ di.Controller = &StartController{}
