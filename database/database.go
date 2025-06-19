package database

import (
	"log"

	"recomemento-api-go/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase initializes the database connection and runs migrations
func InitDatabase(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Book{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}

// SeedDatabase seeds the database with initial data
func SeedDatabase(db *gorm.DB) error {
	// Check if we already have data
	var count int64
	db.Model(&models.Book{}).Count(&count)
	if count > 0 {
		log.Println("Database already contains data, skipping seed")
		return nil
	}

	// Seed data
	books := []models.Book{
		{
			Title:       "The Great Gatsby",
			Author:      "F. Scott Fitzgerald",
			Genre:       "Fiction",
			Purpose:     "Entertainment",
			Description: "A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.",
		},
		{
			Title:       "Clean Code",
			Author:      "Robert C. Martin",
			Genre:       "Technology",
			Purpose:     "Learning",
			Description: "A handbook of agile software craftsmanship that teaches principles of writing clean, readable code.",
		},
		{
			Title:       "1984",
			Author:      "George Orwell",
			Genre:       "Fiction",
			Purpose:     "Entertainment",
			Description: "A dystopian social science fiction novel that follows Winston Smith, a low-ranking citizen of Oceania.",
		},
		{
			Title:       "The Lean Startup",
			Author:      "Eric Ries",
			Genre:       "Business",
			Purpose:     "Learning",
			Description: "A methodology for developing businesses and products that aims to shorten product development cycles.",
		},
	}

	for _, book := range books {
		if err := db.Create(&book).Error; err != nil {
			return err
		}
	}

	log.Printf("Database seeded with %d books", len(books))
	return nil
} 