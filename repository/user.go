package repository

import (
	"bm/entity"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) CreateUser(tgu *gotgbot.User) *entity.User {
	user := &entity.User{
		ID:           tgu.Id,
		FirstName:    tgu.FirstName,
		LastName:     tgu.LastName,
		Username:     tgu.Username,
		LanguageCode: tgu.LanguageCode,
		IsPremium:    tgu.IsPremium,
		State:        entity.IdleState,
	}

	r.db.Create(user)

	return user
}

func (r *UserRepo) FindByID(id int64) *entity.User {
	var u entity.User

	result := r.db.Where("id = ?", id).First(&u)
	if result.Error != nil {
		return nil
	}

	return &u
}
