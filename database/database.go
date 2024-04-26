package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Tax represents tax-related data in the database
type Tax struct {
	gorm.Model
	TotalIncome float64 `gorm:"column:totalIncome"`
	Wht         float64 `gorm:"column:wht"`
	Kreceipt    float64 `gorm:"column:kreceipt"`
	Donation    float64 `gorm:"column:donation"`
	Tax         float64 `gorm:"column:tax"`
}

// CreateTax creates a new tax record in the database
func CreateTax(db *gorm.DB, tax *Tax) {
	result := db.Create(tax)
	if result.Error != nil {
		log.Fatalf("Error creating tax record: %v", result.Error)
	}
	fmt.Println("Tax record created successfully")
}

// GetTax retrieves a tax record from the database by ID
func GetTax(db *gorm.DB, id uint) *Tax {
	var tax Tax
	result := db.First(&tax, id)
	if result.Error != nil {
		log.Fatalf("Error finding tax record: %v", result.Error)
	}
	return &tax
}

// UpdateTax updates an existing tax record in the database
func UpdateTax(db *gorm.DB, tax *Tax) {
	result := db.Save(tax)
	if result.Error != nil {
		log.Fatalf("Error updating tax record: %v", result.Error)
	}
	fmt.Println("Tax record updated successfully")
}

// DeleteTax deletes a tax record from the database by ID
func DeleteTax(db *gorm.DB, id uint) {
	var tax Tax
	result := db.Delete(&tax, id)
	if result.Error != nil {
		log.Fatalf("Error deleting tax record: %v", result.Error)
	}
	fmt.Println("Tax record deleted successfully")
}
