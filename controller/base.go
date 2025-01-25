package controller

import (
	"bm/db"
	"bm/entity"
	"bm/repository"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/fullpipe/icu-mf/mf"
)

// BaseController provides basic shortcuts
type BaseController struct {
	mfBundle mf.Bundle
	userRepo *repository.UserRepo
	em       db.Saver
}

func NewBaseController(
	userRepo *repository.UserRepo,
	mfBundle mf.Bundle,
	em db.Saver,
) BaseController {
	return BaseController{
		mfBundle: mfBundle,
		userRepo: userRepo,
		em:       em,
	}
}

func (c *BaseController) User(ctx *ext.Context) *entity.User {
	user := c.userRepo.FindByID(ctx.EffectiveUser.Id)
	if user == nil {
		user = c.userRepo.CreateUser(ctx.EffectiveUser)
	}

	return user
}

func (c *BaseController) Trans(ctx *ext.Context, id string, args ...mf.TranslationArg) string {
	return c.mfBundle.Translator(ctx.EffectiveUser.LanguageCode).Trans(id, args...)
}

func (c *BaseController) ChangeState(user *entity.User, newState entity.State) error {
	if user.State == newState {
		return nil
	}

	user.State = newState

	return c.em.Save(user, "state")
}
