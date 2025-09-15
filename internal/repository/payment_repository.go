package repository

import (
	"context"

	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (pr *PaymentRepository) Create(ctx context.Context, payment entity.Payment) error {

	err := pr.db.WithContext(ctx).Create(&payment).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil
}

func (pr *PaymentRepository) Update(ctx context.Context, payment entity.Payment) error {

	err := pr.db.WithContext(ctx).Save(&payment).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil
}

func (pr *PaymentRepository) FindByID(ctx context.Context, id string) (entity.Payment, error) {

	query := pr.db.WithContext(ctx).Model(&entity.Payment{})
	var payment entity.Payment

	err := query.First(&payment, id).Error

	if err != nil {
		logrus.Error(err)
		return payment, common.NewError(err, common.ErrResourceNotFound)
	}

	return payment, nil
}

func (pr *PaymentRepository) FindByOrderReference(ctx context.Context, orderReference string) (entity.Payment, error) {

	var payment entity.Payment
	if err := pr.db.Where("order_reference = ?", orderReference).First(&payment).Error; err != nil {
		logrus.Error(err)
		return payment, common.NewError(err, common.ErrResourceNotFound)
	}

	return payment, nil
}
