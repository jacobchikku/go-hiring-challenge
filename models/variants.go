package models

import (
	"github.com/shopspring/decimal"
)

type Variant struct {
	ID        uint   `json:"-" gorm:"primaryKey"`
	ProductID uint   `json:"-" gorm:"not null"`
	Name      string `json:"name" gorm:"not null"`
	SKU       string `json:"sku" gorm:"uniqueIndex;not null"`
	// Note: We use a pointer or check .IsZero() to handle the "optional" price
	Price decimal.Decimal `json:"price" gorm:"type:decimal(10,2);null"`
}

func (v *Variant) TableName() string {
	return "product_variants"
}
