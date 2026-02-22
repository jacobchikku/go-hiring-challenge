package models

type Category struct {
	ID   uint   `json:"-" gorm:"primaryKey"`
	Code string `json:"code" gorm:"uniqueIndex;not null"`
	Name string `json:"name" gorm:"not null"`
}

func (c *Category) TableName() string {
	return "categories"
}
