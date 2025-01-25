package controller

import (
	"bm/di"
	"bm/entity"
	"log/slog"
	"math"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type MessageController struct {
	BaseController
}

func NewMessageController(
	bc BaseController,
) *MessageController {
	return &MessageController{
		BaseController: bc,
	}
}

func (c *MessageController) Message(b *gotgbot.Bot, ctx *ext.Context) error {
	user := c.User(ctx)
	switch user.State {

	default:
		slog.Info(
			"Unable to handle message",
			slog.Int64("user_id", ctx.EffectiveUser.Id),
			slog.String("message", ctx.Message.Text),
		)

		return c.ChangeState(user, entity.IdleState)
	}
}

func (c *MessageController) Register(bot *gotgbot.Bot, dispatcher *ext.Dispatcher) error {
	dispatcher.AddHandler(handlers.NewMessage(message.All, c.Message))

	return nil
}

func (c *MessageController) Priority() int {
	return math.MinInt
}

var _ di.ControllerWithPriority = &MessageController{}
