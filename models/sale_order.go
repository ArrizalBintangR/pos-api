package models

import (
	"time"

	"gorm.io/gorm"
)

type SaleOrder struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	OrderNumber    string          `gorm:"uniqueIndex;not null;size:50" json:"order_number"`
	CustomerName   string          `gorm:"size:255;not null" json:"customer_name"`
	TotalAmount    float64         `gorm:"not null;default:0" json:"total_amount"`
	Notes          string          `gorm:"type:text" json:"notes"`
	CreatedByID    uint            `gorm:"not null" json:"created_by_id"`
	CreatedBy      *User           `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	SaleOrderItems []SaleOrderItem `gorm:"foreignKey:SaleOrderID" json:"items,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (SaleOrder) TableName() string {
	return "sale_orders"
}

type SaleOrderItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SaleOrderID uint           `gorm:"not null" json:"sale_order_id"`
	ProductName string         `gorm:"not null;size:255" json:"product_name"`
	Quantity    int            `gorm:"not null;default:1" json:"quantity"`
	UnitPrice   float64        `gorm:"not null" json:"unit_price"`
	Subtotal    float64        `gorm:"not null" json:"subtotal"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SaleOrderItem) TableName() string {
	return "sale_order_items"
}
