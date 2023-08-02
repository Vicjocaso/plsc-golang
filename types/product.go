package types

// A Product contains metadata about a product for sale.
type Product struct {
	// gorm.Model
	ID          int
	Name        string
	Description string
	Image       string
	Price       float64
	Inventory   int
	Tags        string
	CategoryID  int
	Category    Category `gorm:"foreignKey:CategoryID"`
}

// A Category describes a group of Products.
type Category struct {
	ID          int
	Name        string
	Description string
}
