package di

import (
	"bm/config"
	"bm/i18n"
	"context"
	"log/slog"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/fullpipe/icu-mf/mf"
	"go.uber.org/fx"
)

func InitBot(lc fx.Lifecycle, cfg config.Config, mfBundle mf.Bundle) (*gotgbot.Bot, *ext.Dispatcher, error) {
	b, err := gotgbot.NewBot(cfg.Token, nil)
	if err != nil {
		return nil, nil, err
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			slog.Error(err.Error(), slog.Int64("tg_id", ctx.EffectiveUser.Id))

			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	lc.Append(fx.StartHook(func(ctx context.Context) {
		i18n.LocalizeBot(b, mfBundle)
	}))

	lc.Append(fx.StopHook(func(ctx context.Context) {
	}))

	return b, dispatcher, nil
}

func InitBotUpdater(
	lc fx.Lifecycle,
	bot *gotgbot.Bot,
	dispatcher *ext.Dispatcher,
	controllers ...Controller,
) (*ext.Updater, error) {
	updater := ext.NewUpdater(dispatcher, nil)

	err := updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: false,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		for _, controller := range controllers {
			err := controller.Register(bot, dispatcher)
			if err != nil {
				return err
			}
		}

		return nil
	}))

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		return updater.Stop()
	}))

	return updater, nil
}
