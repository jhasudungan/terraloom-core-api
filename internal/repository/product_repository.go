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

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) FindWithFilters(
	ctx context.Context,
	filter model.ProductFilter,
	pagination model.PaginationParams) ([]entity.Product, int64, error) {

	baseQuery := pr.db.WithContext(ctx).Model(&entity.Product{})
	var products []entity.Product
	var total int64

	// Apply filters
	if filter.Name != "" {
		baseQuery = baseQuery.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.IsActive {
		baseQuery = baseQuery.Where("is_active = ?", filter.IsActive)
	}

	// Exclude soft deleted
	baseQuery = baseQuery.Where("deleted_at IS NULL")

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
	if err := dataQuery.Find(&products).Error; err != nil {
		logrus.Error(err)
		return nil, 0, common.NewError(err, common.ErrResourceNotFound)
	}

	return products, total, nil
}

func (pr *ProductRepository) FindByID(ctx context.Context, id int64) (entity.Product, error) {

	query := pr.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"})
	var product entity.Product

	err := query.First(&product, id).Error

	if err != nil {
		return product, common.NewError(err, common.ErrResourceNotFound)
	}

	return product, nil
}

func (pr *ProductRepository) FindMultipleByIDs(ctx context.Context, ids []int64) ([]entity.Product, error) {

	var products []entity.Product
	err := pr.db.WithContext(ctx).Where("id IN ?", ids).Find(&products).Error

	if err != nil {
		return products, common.NewError(err, common.ErrDBOperation)
	}

	return products, nil
}

func (pr *ProductRepository) CheckById(ctx context.Context, id int64) (bool, error) {

	var count int64

	err := pr.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Count(&count).Error

	if err != nil {
		logrus.Error(err)
		return false, common.NewError(err, common.ErrDBOperation)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (pr *ProductRepository) Update(ctx context.Context, product entity.Product) error {

	err := pr.db.Save(&product).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil
}

func (pr *ProductRepository) BatchUpsert(ctx context.Context, products []entity.Product) error {

	err := pr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"stock", "updated_at", "updated_by"}),
	}).Create(&products).Error

	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrDBOperation)
	}

	return nil
}
