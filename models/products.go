package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID         uint            `json:"-" gorm:"primaryKey"`
	Code       string          `json:"code" gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not null"`
	CategoryID uint            `json:"-"`
	// Task 2 & 3: Relationship to Category
	Category Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Variants []Variant `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
}

func (p *Product) TableName() string {
	return "products"
}

// AfterFind handles Requirement: "variants without specific price should inherit the price from the product"
func (p *Product) AfterFind(tx *gorm.DB) (err error) {
	for i := range p.Variants {
		// If the variant price is 0 (not set in DB), use the parent product's price
		if p.Variants[i].Price.IsZero() {
			p.Variants[i].Price = p.Price
		}
	}
	return nil
}
