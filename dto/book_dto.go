package dto

// CreateBookRequest represents the request body for creating a book
type CreateBookRequest struct {
	// The title of the book
	Title string `json:"title" binding:"required" example:"The Great Gatsby"`
	// The author of the book
	Author string `json:"author" binding:"required" example:"F. Scott Fitzgerald"`
	// The genre of the book
	Genre string `json:"genre" binding:"required" example:"Fiction"`
	// The purpose of the book
	Purpose string `json:"purpose" binding:"required" example:"Entertainment"`
	// The description of the book
	Description string `json:"description" binding:"required" example:"A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan."`
}

// UpdateBookRequest represents the request body for updating a book
// All fields are optional and will only update the provided fields
type UpdateBookRequest struct {
	// The title of the book (optional)
	Title *string `json:"title,omitempty" example:"The Great Gatsby"`
	// The author of the book (optional)
	Author *string `json:"author,omitempty" example:"F. Scott Fitzgerald"`
	// The genre of the book (optional)
	Genre *string `json:"genre,omitempty" example:"Fiction"`
	// The purpose of the book (optional)
	Purpose *string `json:"purpose,omitempty" example:"Entertainment"`
	// The description of the book (optional)
	Description *string `json:"description,omitempty" example:"A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan."`
}

// RecommendBookRequest represents the request body for book recommendation
type RecommendBookRequest struct {
	// The genre to search for recommendations
	Genre string `json:"genre" binding:"required" example:"Fiction"`
	// The type of book (optional)
	Type string `json:"type" example:"Novel"`
	// The purpose of the book for recommendation
	Purpose string `json:"purpose" binding:"required" example:"Entertainment"`
}

// BookResponse represents the response body for book operations
type BookResponse struct {
	// Unique identifier for the book
	ID uint `json:"id" example:"1"`
	// The title of the book
	Title string `json:"title" example:"The Great Gatsby"`
	// The author of the book
	Author string `json:"author" example:"F. Scott Fitzgerald"`
	// The genre of the book
	Genre string `json:"genre" example:"Fiction"`
	// The purpose of the book
	Purpose string `json:"purpose" example:"Entertainment"`
	// The description of the book
	Description string `json:"description" example:"A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan."`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	// Error type or code
	Error string `json:"error" example:"Book not found"`
	// Detailed error message
	Message string `json:"message" example:"The requested book could not be found"`
} 