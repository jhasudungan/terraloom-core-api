package repository

import (
	"context"

	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (oir *OrderItemRepository) Create(ctx context.Context, orderItem entity.OrderItem) error {

	err := oir.db.WithContext(ctx).Create(&orderItem).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil
}

func (oir *OrderItemRepository) FindByID(ctx context.Context, id string) (entity.OrderItem, error) {

	query := oir.db.WithContext(ctx).Model(&entity.OrderItem{})
	var orderItem entity.OrderItem

	err := query.Where("order_item_reference = ?", id).First(&orderItem).Error

	if err != nil {
		logrus.Error(err)
		return orderItem, common.NewError(err, common.ErrResourceNotFound)
	}

	return orderItem, nil

}

func (oir *OrderItemRepository) FindByOrderID(ctx context.Context, orderID uint) ([]entity.OrderItem, error) {

	var orderItems []entity.OrderItem

	err := oir.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&orderItems).Error

	if err != nil {
		logrus.Error(err)
		return orderItems, common.NewError(err, common.ErrResourceNotFound)
	}

	return orderItems, nil
}

func (oir *OrderItemRepository) CreateBatch(ctx context.Context, orderItems []entity.OrderItem, batchSize int) error {

	err := oir.db.WithContext(ctx).CreateInBatches(&orderItems, batchSize).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil
}
