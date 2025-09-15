package entity

import "time"

type Order struct {
	OrderReference  string     `gorm:"primaryKey;column:order_reference"`
	OrderDate       time.Time  `gorm:"column:order_date;default:CURRENT_TIMESTAMP"`
	AccountUsername string     `gorm:"column:account_username"`
	DeliveryAddress string     `gorm:"column:delivery_address"`
	Status          string     `gorm:"column:status;default:PENDING;size:200"`
	Total           int64      `gorm:"column:total;default:0"`
	CreatedAt       time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy       string     `gorm:"column:created_by;size:100"`
	UpdatedBy       string     `gorm:"column:updated_by;size:100"`
	DeletedAt       *time.Time `gorm:"column:deleted_at"`

	// Relationship: One order has many order items
	OrderItems []OrderItem `gorm:"foreignKey:OrderReference;references:OrderReference"`

	// Pointer avoids recursive allocation
	Payment *Payment `gorm:"foreignKey:OrderReference;references:OrderReference"`

	Account Account `gorm:"foreignKey:AccountUsername;references:Username"`
}

func (Order) TableName() string {
	return "orders"
}
