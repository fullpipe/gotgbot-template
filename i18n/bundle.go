package i18n

import (
	"embed"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/fullpipe/icu-mf/mf"
	"golang.org/x/text/language"
)

//go:embed translations/messages.*.yaml
var messagesDir embed.FS

func NewMessageBundle() (mf.Bundle, error) {
	return mf.NewBundle(
		mf.WithDefaulLangFallback(language.English),
		mf.WithErrorHandler(func(err error, id string, ctx map[string]any) {
			slog.Error(err.Error(), slog.String("id", id), slog.Any("ctx", ctx))
		}),
		mf.WithYamlProvider(messagesDir),
	)
}

func LocalizeBot(b *gotgbot.Bot, bundle mf.Bundle) {
	en := bundle.Translator("en")
	ru := bundle.Translator("ru")

	b.SetMyName(&gotgbot.SetMyNameOpts{
		Name: en.Trans("bot.name"),
	})
	b.SetMyName(&gotgbot.SetMyNameOpts{
		Name:         ru.Trans("bot.name"),
		LanguageCode: "ru",
	})

	b.SetMyDescription(&gotgbot.SetMyDescriptionOpts{
		Description: en.Trans("bot.description"),
	})
	b.SetMyDescription(&gotgbot.SetMyDescriptionOpts{
		Description:  ru.Trans("bot.description"),
		LanguageCode: "ru",
	})

	b.SetMyShortDescription(&gotgbot.SetMyShortDescriptionOpts{
		ShortDescription: en.Trans("bot.short_description"),
	})
	b.SetMyShortDescription(&gotgbot.SetMyShortDescriptionOpts{
		ShortDescription: ru.Trans("bot.short_description"),
		LanguageCode:     "ru",
	})
}
