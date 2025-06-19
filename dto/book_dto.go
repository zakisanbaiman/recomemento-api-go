package dto

// CreateBookRequest represents the request body for creating a book
type CreateBookRequest struct {
	Title       string `json:"title" binding:"required" example:"The Great Gatsby"`
	Author      string `json:"author" binding:"required" example:"F. Scott Fitzgerald"`
	Genre       string `json:"genre" binding:"required" example:"Fiction"`
	Purpose     string `json:"purpose" binding:"required" example:"Entertainment"`
	Description string `json:"description" binding:"required" example:"A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan."`
}

// UpdateBookRequest represents the request body for updating a book
type UpdateBookRequest struct {
	Title       *string `json:"title,omitempty" example:"The Great Gatsby"`
	Author      *string `json:"author,omitempty" example:"F. Scott Fitzgerald"`
	Genre       *string `json:"genre,omitempty" example:"Fiction"`
	Purpose     *string `json:"purpose,omitempty" example:"Entertainment"`
	Description *string `json:"description,omitempty" example:"A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan."`
}

// RecommendBookRequest represents the request body for book recommendation
type RecommendBookRequest struct {
	Genre   string `json:"genre" binding:"required" example:"Fiction"`
	Type    string `json:"type" example:"Novel"`
	Purpose string `json:"purpose" binding:"required" example:"Entertainment"`
}

// BookResponse represents the response body for book operations
type BookResponse struct {
	ID          uint   `json:"id" example:"1"`
	Title       string `json:"title" example:"The Great Gatsby"`
	Author      string `json:"author" example:"F. Scott Fitzgerald"`
	Genre       string `json:"genre" example:"Fiction"`
	Purpose     string `json:"purpose" example:"Entertainment"`
	Description string `json:"description" example:"A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan."`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Book not found"`
	Message string `json:"message" example:"The requested book could not be found"`
} 