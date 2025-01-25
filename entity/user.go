package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Greeter interface {
	Greet() string
}

type User struct {
	ID int64 `gorm:"primarykey"`

	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
	IsPremium    bool

	State State

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type State string

const (
	IdleState State = "idle"
)

func (u *User) Fullname() string {
	return fmt.Sprintf("%s %s (@%s)", u.FirstName, u.LastName, u.Username)
}

func (u *User) TgURL() string {
	return fmt.Sprintf("tg://user?id=%d", u.ID)
}
