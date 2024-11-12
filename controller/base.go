package controller

import (
	"bm/entity"
	"bm/repository"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/fullpipe/icu-mf/mf"
)

// baseController provides basic shortcuts
type baseController struct {
	mfBundle mf.Bundle
	userRepo *repository.UserRepo
}

func (c *baseController) User(ctx *ext.Context) *entity.User {
	user := c.userRepo.FindByID(ctx.EffectiveUser.Id)
	if user == nil {
		user = c.userRepo.CreateUser(ctx.EffectiveUser)
	}

	return user
}

func (c *baseController) Trans(ctx *ext.Context, id string, args ...mf.TranslationArg) string {
	return c.mfBundle.Translator(ctx.EffectiveUser.LanguageCode).Trans(id, args...)
}
