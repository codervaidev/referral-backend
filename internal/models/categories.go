package models

// Category represents a product category (categories table).
type Category struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description *string `json:"description"`
}

func (Category) TableName() string { return "categories" } 