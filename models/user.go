package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleCashier Role = "cashier"
	RoleOwner   Role = "owner"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	Name      string         `gorm:"not null;size:255" json:"name"`
	Role      Role           `gorm:"not null;size:20" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
