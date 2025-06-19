package models

import "gorm.io/gorm"

// Book represents the book model
type Book struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title" gorm:"not null" binding:"required"`
	Author      string `json:"author" gorm:"not null" binding:"required"`
	Genre       string `json:"genre" gorm:"not null" binding:"required"`
	Purpose     string `json:"purpose" gorm:"not null" binding:"required"`
	Description string `json:"description" gorm:"not null" binding:"required"`
}

// TableName specifies the table name for the Book model
func (Book) TableName() string {
	return "books"
}

// BookDatabase interface for book operations
type BookDatabase interface {
	Create(book *Book) error
	GetAll() ([]Book, error)
	GetByID(id uint) (*Book, error)
	Update(id uint, updates map[string]interface{}) (*Book, error)
	Delete(id uint) (*Book, error)
	FindByGenreAndPurpose(genre, purpose string) (*Book, error)
}

// bookRepository implements BookDatabase
type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository creates a new book repository
func NewBookRepository(db *gorm.DB) BookDatabase {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) GetAll() ([]Book, error) {
	var books []Book
	err := r.db.Find(&books).Error
	return books, err
}

func (r *bookRepository) GetByID(id uint) (*Book, error) {
	var book Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) Update(id uint, updates map[string]interface{}) (*Book, error) {
	var book Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	
	err = r.db.Model(&book).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	
	return &book, nil
}

func (r *bookRepository) Delete(id uint) (*Book, error) {
	var book Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	
	err = r.db.Delete(&book).Error
	if err != nil {
		return nil, err
	}
	
	return &book, nil
}

func (r *bookRepository) FindByGenreAndPurpose(genre, purpose string) (*Book, error) {
	var book Book
	err := r.db.Where("genre = ? AND purpose = ?", genre, purpose).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
} 