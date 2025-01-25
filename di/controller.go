package di

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/fx"
)

func AsController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Controller)),
		fx.ResultTags(`group:"controllers"`),
	)
}

type Controller interface {
	// Register registeres controller handlers
	Register(bot *gotgbot.Bot, dispatcher *ext.Dispatcher) error
}

type ControllerWithPriority interface {
	Controller

	// Priority returns controller priority, default 0
	Priority() int
}
