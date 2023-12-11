package main

import (
	"errors"
	"fmt"
	"github.com/grahms/xgor"
	"gorm.io/driver/sqlite"
)

// Product Define a sample entity
type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Price float64
}

func main() {
	// Connect to an in-memory SQLite database
	db, err := xgor.Open(sqlite.Open("file::memory:"))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	// Migrate the database schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		fmt.Println("Error migrating database schema:", err)
		return
	}

	// Create a repository for the Product entity
	productRepo := xgor.New[Product](db, errors.New("product not found"))

	// Add a product
	newProduct := &Product{Name: "Laptop", Price: 999.99}
	err = productRepo.Add(newProduct)
	if err != nil {
		fmt.Println("Error adding product:", err)
		return
	}

	// Update the product
	updatedProduct := &Product{ID: newProduct.ID, Name: "Updated Laptop", Price: 1099.99}
	err = productRepo.Update(updatedProduct)
	if err != nil {
		fmt.Println("Error updating product:", err)
		return
	}

	// Retrieve a product by ID
	retrievedProduct, err := productRepo.GetByID(newProduct.ID)
	if err != nil {
		fmt.Println("Error getting product by ID:", err)
		return
	}
	fmt.Println("Retrieved Product:", retrievedProduct)

	// Use custom filters to get products with a specific condition
	filters := xgor.FilterType{"price__gt": 1000.0}
	highPricedProducts, err := productRepo.GetAll(nil, nil, nil, filters)
	if err != nil {
		fmt.Println("Error getting high-priced products:", err)
		return
	}
	fmt.Println("High-Priced Products:", highPricedProducts.Items)
}
