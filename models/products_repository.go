package models

import (
	"gorm.io/gorm"
)

// ProductRepository defines the behavior for product and category data access
type ProductRepository interface {
	GetProducts(limit, offset int, category string, priceLessThan float64) ([]Product, int64, error)
	GetByCode(code string) (*Product, error)
	GetAllCategories() ([]Category, error)
	CreateCategory(category *Category) error
}

type productsRepository struct {
	db *gorm.DB
}

// NewProductsRepository creates a new instance of the repository
func NewProductsRepository(db *gorm.DB) ProductRepository {
	return &productsRepository{
		db: db,
	}
}

func (r *productsRepository) GetProducts(limit, offset int, category string, priceLessThan float64) ([]Product, int64, error) {
	var products []Product
	var total int64

	// Initialize query with Preloads for Category and Variants
	query := r.db.Model(&Product{}).Preload("Category").Preload("Variants")

	// Filter by Category Code
	if category != "" {
		query = query.Joins("Category").Where("category.code = ?", category)
	}

	// Filter by Price
	if priceLessThan > 0 {
		query = query.Where("products.price < ?", priceLessThan)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	err := query.Limit(limit).Offset(offset).Find(&products).Error

	return products, total, err
}

func (r *productsRepository) GetByCode(code string) (*Product, error) {
	var product Product
	err := r.db.Preload("Category").Preload("Variants").
		Where("code = ?", code).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productsRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *productsRepository) CreateCategory(category *Category) error {
	return r.db.Create(category).Error
}
