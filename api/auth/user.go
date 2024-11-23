package auth

import (
	"bm/entity"
	"context"

	initdata "github.com/telegram-mini-apps/init-data-golang"
	"gorm.io/gorm"
)

type userContext struct {
	db       *gorm.DB
	initData initdata.InitData
	user     *entity.User
}

func User(ctx context.Context) *entity.User {
	uc, ok := ctx.Value(userContextKey).(*userContext)
	if !ok {
		return nil
	}

	if uc.user != nil {
		return uc.user
	}

	if uc.initData.User.ID == 0 {
		return nil
	}

	var user entity.User
	result := uc.db.First(&user, uc.initData.User.ID)
	if result.Error != nil {
		return nil
	}

	uc.user = &user

	return uc.user
}

func UserID(ctx context.Context) int64 {
	raw, ok := ctx.Value(userContextKey).(*userContext)
	if !ok {
		return 0
	}

	return raw.initData.User.ID
}

func Roles(ctx context.Context) []string {
	_, ok := ctx.Value(userContextKey).(*userContext)
	if !ok {
		return []string{}
	}

	// TODO: roles from user model
	return []string{"USER"}
}
