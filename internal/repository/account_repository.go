package repository

import (
	"context"

	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (ar *AccountRepository) Create(ctx context.Context, account entity.Account) error {

	err := ar.db.WithContext(ctx).Create(&account).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil
}

func (ar *AccountRepository) Update(ctx context.Context, account entity.Account) error {

	err := ar.db.WithContext(ctx).Save(&account).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil

}

func (ar *AccountRepository) FindByUsername(ctx context.Context, username string) (entity.Account, error) {

	query := ar.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"})
	var account entity.Account

	err := query.Where("username = ?", username).First(&account).Error

	if err != nil {
		return account, common.NewError(err, common.ErrResourceNotFound)
	}

	return account, nil
}

func (ar *AccountRepository) CheckByUsername(ctx context.Context, username string) (bool, error) {

	var count int64

	err := ar.db.WithContext(ctx).Model(&entity.Account{}).Where("username = ?", username).Count(&count).Error

	if err != nil {
		logrus.Error(err)
		return false, common.NewError(err, common.ErrDBOperation)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (ar *AccountRepository) CheckByEmail(ctx context.Context, email string) (bool, error) {

	var count int64

	err := ar.db.WithContext(ctx).Model(&entity.Account{}).Where("email = ?", email).Count(&count).Error

	if err != nil {
		logrus.Error(err)
		return false, common.NewError(err, common.ErrDBOperation)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
