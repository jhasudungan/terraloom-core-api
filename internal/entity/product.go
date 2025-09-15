package entity

import "time"

type Product struct {
	ID          int64      `gorm:"primaryKey;column:id"`
	CategoryID  int64      `gorm:"column:category_id"`
	Name        string     `gorm:"column:name"`
	Description string     `gorm:"column:description"`
	Stock       int64      `gorm:"column:stock"`
	Price       int64      `gorm:"column:price"`
	ImageUrl    string     `gorm:"column:image_url"`
	IsActive    bool       `gorm:"column:is_active"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	CreatedBy   string     `gorm:"column:created_by"`
	UpdatedBy   string     `gorm:"column:updated_by"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (Product) TableName() string {
	return "products"
}

func (p *Product) IsDeleted() bool {
	return p.DeletedAt != nil
}
