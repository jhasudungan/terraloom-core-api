package entity

import "time"

type Payment struct {
	PaymentReference string     `gorm:"primaryKey;column:payment_reference"`
	OrderReference   string     `gorm:"not null;index"`
	Total            int64      `gorm:"column:total"`
	CardHolderName   string     `gorm:"column:card_holder_name"`
	CardNumber       string     `gorm:"column:card_number"`
	Status           string     `gorm:"column:status"`
	PaymentDate      time.Time  `gorm:"column:payment_date;default:CURRENT_TIMESTAMP"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
	CreatedBy        string     `gorm:"column:created_by"`
	UpdatedBy        string     `gorm:"column:updated_by"`
	DeletedAt        *time.Time `gorm:"column:deleted_at"`

	Order Order `gorm:"foreignKey:OrderReference;references:OrderReference"`
}

func (Payment) TableName() string {
	return "payments"
}
