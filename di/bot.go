package di

import (
	"bm/config"
	"bm/i18n"
	"bm/tools"
	"cmp"
	"context"
	"log/slog"
	"slices"
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

	en := mfBundle.Translator("en")
	ru := mfBundle.Translator("ru")
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.StartHook(func() {
		i18n.LocalizeBot(b, mfBundle)
		tools.Retry(ctx, func(ctx context.Context) (any, error) {
			return b.SetMyCommandsWithContext(ctx, []gotgbot.BotCommand{
				{Command: "support", Description: en.Trans("bot.commands.support")},
			}, nil)
		})

		tools.Retry(ctx, func(ctx context.Context) (any, error) {
			return b.SetMyCommandsWithContext(ctx, []gotgbot.BotCommand{
				{Command: "support", Description: ru.Trans("bot.commands.support")},
			}, &gotgbot.SetMyCommandsOpts{
				Scope:        gotgbot.BotCommandScopeDefault{},
				LanguageCode: "ru",
			})
		})

		tools.Retry(ctx, func(ctx context.Context) (any, error) {
			return b.SetMyNameWithContext(ctx, &gotgbot.SetMyNameOpts{
				Name: en.Trans("bot.name"),
			})
		})

		tools.Retry(ctx, func(ctx context.Context) (any, error) {
			return b.SetMyNameWithContext(ctx, &gotgbot.SetMyNameOpts{
				Name:         ru.Trans("bot.name"),
				LanguageCode: "ru",
			})
		})

		tools.Retry(ctx, func(ctx context.Context) (any, error) {
			return b.SetMyDescriptionWithContext(ctx, &gotgbot.SetMyDescriptionOpts{
				Description: en.Trans("bot.description"),
			})
		})

		tools.Retry(ctx, func(ctx context.Context) (any, error) {
			return b.SetMyDescriptionWithContext(ctx, &gotgbot.SetMyDescriptionOpts{
				Description:  ru.Trans("bot.description"),
				LanguageCode: "ru",
			})
		})

		tools.Retry(ctx, func(ctx context.Context) (any, error) {
			return b.SetMyShortDescriptionWithContext(ctx, &gotgbot.SetMyShortDescriptionOpts{
				ShortDescription: en.Trans("bot.short_description"),
			})
		})

		tools.Retry(ctx, func(ctx context.Context) (any, error) {
			return b.SetMyShortDescriptionWithContext(ctx, &gotgbot.SetMyShortDescriptionOpts{
				ShortDescription: ru.Trans("bot.short_description"),
				LanguageCode:     "ru",
			})
		})
	}))

	lc.Append(fx.StopHook(func() {
		cancel()
	}))

	return b, dispatcher, nil
}

func InitBotUpdater(
	lc fx.Lifecycle,
	bot *gotgbot.Bot,
	dispatcher *ext.Dispatcher,
	controllers ...Controller,
) (*ext.Updater, error) {
	slices.SortFunc(controllers, func(a, b Controller) int {
		pa, pb := 0, 0
		if a, ok := a.(ControllerWithPriority); ok {
			pa = a.Priority()
		}

		if b, ok := b.(ControllerWithPriority); ok {
			pb = b.Priority()
		}

		return cmp.Compare(pb, pa)
	})

	for _, controller := range controllers {
		err := controller.Register(bot, dispatcher)
		if err != nil {
			return nil, err
		}
	}

	updater := ext.NewUpdater(dispatcher, nil)

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		// We start polling after registering controllers
		return updater.StartPolling(bot, &ext.PollingOpts{
			DropPendingUpdates: false,
			GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
				Timeout: 9,
				RequestOpts: &gotgbot.RequestOpts{
					Timeout: time.Second * 10,
				},
			},
		})
	}))

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		return updater.Stop()
	}))

	return updater, nil
}
