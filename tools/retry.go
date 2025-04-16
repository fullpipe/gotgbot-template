package tools

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func Retry[T any](ctx context.Context, f func(ctx context.Context) (T, error)) {
	go func() {
		_, err := f(ctx)

		for {
			tgErr := &gotgbot.TelegramError{}
			if errors.As(err, &tgErr) && tgErr.ResponseParams != nil && tgErr.ResponseParams.RetryAfter > 0 {
				slog.Warn(
					"Init retry: "+tgErr.Error(),
					slog.Any("RetryAfter", tgErr.ResponseParams.RetryAfter),
					slog.Any("err", *tgErr),
				)

				select {
				case <-ctx.Done():
					slog.Warn("retry context canceled")
					return
				case <-time.After(time.Duration(tgErr.ResponseParams.RetryAfter+rand.Int64N(200)) * time.Second):
					slog.Debug(
						"Execute retry: "+tgErr.Error(),
						slog.Any("RetryAfter", tgErr.ResponseParams.RetryAfter),
					)
					_, err = f(ctx)
				}
			} else {
				if err != nil {
					slog.Error(err.Error())
				}

				return
			}
		}
	}()
}
