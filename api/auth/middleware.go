package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
	"gorm.io/gorm"
)

var userContextKey contextKey

type contextKey int

// Middleware which authorizes the external client.
func Middleware(token string, db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// We expect passing init data in the Authorization header in the following format:
			// <auth-type> <auth-data>
			// <auth-type> must be "tma", and <auth-data> is Telegram Mini Apps init data.
			authParts := strings.Fields(r.Header.Get("authorization"))
			if len(authParts) != 2 {
				// TODO: all users should be authed??
				http.Error(w, "Auth token is required", http.StatusUnauthorized)
				// next.ServeHTTP(w, r)
				return
			}

			authType := authParts[0]
			authData := authParts[1]

			slog.Info("authData", slog.Any("authData", authData))

			switch authType {
			case "tma":
				// Validate init data. We consider init data sign valid for 1 hour from their
				// creation moment.
				if err := initdata.Validate(authData, token, time.Hour); err != nil {
					http.Error(w, "Invalid init data", http.StatusUnauthorized)
					return
				}

				// Parse init data. We will surely need it in the future.
				initData, err := initdata.Parse(authData)
				if err != nil {
					http.Error(w, "Unable to parse init data", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), userContextKey, &userContext{
					db:       db,
					initData: initData,
				})

				// and call the next with our new context
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			}
		})
	}
}
