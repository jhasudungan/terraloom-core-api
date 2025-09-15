package entity

import "time"

type OrderItem struct {
	OrderItemReference      string     `gorm:"primaryKey;column:order_item_reference"`
	OrderReference          string     `gorm:"not null;index"`
	ProductID               int64      `gorm:"column:product_id"`
	PriceSnapshot           int64      `gorm:"column:price_snapshot;default:0"`
	Quantity                int64      `gorm:"column:quantity;default:0"`
	Total                   int64      `gorm:"column:total;default:0"`
	ProductNameSnapshot     string     `gorm:"column:product_name_snapshot"`
	ProductImageUrlSnapshot string     `gorm:"column:product_image_url_snapshot"`
	CreatedAt               time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt               time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy               string     `gorm:"column:created_by;size:100"`
	UpdatedBy               string     `gorm:"column:updated_by;size:100"`
	DeletedAt               *time.Time `gorm:"column:deleted_at"`

	// Relationship: Each order item belongs to one order
	Order Order `gorm:"foreignKey:OrderReference;references:OrderReference"`

	Product Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
