package entity

import "time"

type Account struct {
	ID                int64      `gorm:"primaryKey;column:id"`
	Username          string     `gorm:"column:username"`
	DisplayName       string     `gorm:"column:display_name"`
	Email             string     `gorm:"column:email"`
	LoginPassword     string     `gorm:"column:login_password"`
	RegisteredAddress string     `gorm:"column:registered_address"`
	IsActive          bool       `gorm:"column:is_active"`
	CreatedAt         time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy         string     `gorm:"column:created_by;size:100"`
	UpdatedBy         string     `gorm:"column:updated_by;size:100"`
	DeletedAt         *time.Time `gorm:"column:deleted_at"`

	Orders []Order `gorm:"foreignKey:AccountUsername;references:Username"`
}

func (Account) TableName() string {
	return "accounts"
}
