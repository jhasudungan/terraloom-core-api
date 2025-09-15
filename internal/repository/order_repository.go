package repository

import (
	"context"

	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/entity"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) Create(ctx context.Context, order entity.Order) error {

	err := or.db.WithContext(ctx).Create(&order).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil
}

func (or *OrderRepository) Update(ctx context.Context, order entity.Order) error {

	err := or.db.WithContext(ctx).Save(&order).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil

}

func (or *OrderRepository) FindByID(ctx context.Context, id string) (entity.Order, error) {

	query := or.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"})
	var order entity.Order

	err := query.Where("order_reference = ?", id).First(&order).Error

	if err != nil {
		return order, common.NewError(err, common.ErrResourceNotFound)
	}

	return order, nil
}

func (or *OrderRepository) FindByIDWithItems(ctx context.Context, id string) (entity.Order, error) {

	query := or.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"})
	var order entity.Order

	err := query.Preload("OrderItems").Where("order_reference = ?", id).First(&order).Error

	if err != nil {
		return order, common.NewError(err, common.ErrResourceNotFound)
	}

	return order, nil
}

func (or *OrderRepository) FindWithAccountAndFilters(
	ctx context.Context,
	accountUsername string,
	filter model.OrderFilter,
	pagination model.PaginationParams) ([]entity.Order, int64, error) {

	baseQuery := or.db.WithContext(ctx).Model(&entity.Order{})
	var orders []entity.Order
	var total int64

	// Apply filters
	if filter.OrderReference != "" {
		baseQuery = baseQuery.Where("order_reference ILIKE ?", "%"+filter.OrderReference+"%")
	}

	// account_username
	baseQuery = baseQuery.Where("account_username", accountUsername)

	// Exclude soft deleted
	baseQuery = baseQuery.Where("deleted_at IS NULL")

	// Add ordering by order_date (newest first)
	baseQuery = baseQuery.Order("order_date DESC")

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		logrus.Error(err)
		return nil, 0, common.NewError(err, common.ErrResourceNotFound)
	}

	// Apply pagination if enabled
	dataQuery := baseQuery
	if pagination.IsPaginate {
		dataQuery = dataQuery.Offset(pagination.GetOffset()).Limit(pagination.PerPage)
	}

	// Execute query
	if err := dataQuery.Find(&orders).Error; err != nil {
		logrus.Error(err)
		return nil, 0, common.NewError(err, common.ErrResourceNotFound)
	}

	return orders, total, nil
}

func (or *OrderRepository) GetDB() *gorm.DB {
	return or.db
}
