package db

import "gorm.io/gorm"

type Saver interface {
	Save(entity any, fields ...string) error
}

type Em struct {
	db *gorm.DB
}

func NewEm(db *gorm.DB) *Em {
	return &Em{
		db: db,
	}
}

func (em *Em) Save(entity any, fields ...string) error {
	if len(fields) == 0 {
		return em.db.Save(entity).Error
	}

	return em.db.Model(entity).Select(fields).Updates(entity).Error
}
