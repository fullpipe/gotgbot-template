package repository

import (
	"bm/entity"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db}
}

func (r *PaymentRepo) FindByID(id uint) *entity.Payment {
	var p entity.Payment

	result := r.db.Where("id = ?", id).First(&p)
	if result.Error != nil {
		return nil
	}

	return &p
}
