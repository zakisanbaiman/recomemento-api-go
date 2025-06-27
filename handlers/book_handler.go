package handlers

import (
	"net/http"
	"strconv"

	"recomemento-api-go/dto"
	"recomemento-api-go/models"

	"github.com/gin-gonic/gin"
)

// BookHandler handles book-related HTTP requests
type BookHandler struct {
	bookRepo models.BookDatabase
}

// NewBookHandler creates a new book handler
func NewBookHandler(bookRepo models.BookDatabase) *BookHandler {
	return &BookHandler{
		bookRepo: bookRepo,
	}
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with the provided information
// @Tags books
// @Accept json
// @Produce json
// @Param book body dto.CreateBookRequest true "Book information"
// @Success 201 {object} dto.BookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req dto.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	book := &models.Book{
		Title:       req.Title,
		Author:      req.Author,
		Genre:       req.Genre,
		Purpose:     req.Purpose,
		Description: req.Description,
	}

	if err := h.bookRepo.Create(book); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create book",
			Message: err.Error(),
		})
		return
	}

	response := dto.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Genre:       book.Genre,
		Purpose:     book.Purpose,
		Description: book.Description,
	}

	c.JSON(http.StatusCreated, response)
}

// GetAllBooks godoc
// @Summary Get all books
// @Description Get a list of all books
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} dto.BookResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books [get]
func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.bookRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get books",
			Message: err.Error(),
		})
		return
	}

	var response []dto.BookResponse
	for _, book := range books {
		response = append(response, dto.BookResponse{
			ID:          book.ID,
			Title:       book.Title,
			Author:      book.Author,
			Genre:       book.Genre,
			Purpose:     book.Purpose,
			Description: book.Description,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetBookByID godoc
// @Summary Get a book by ID
// @Description Get a specific book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} dto.BookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
		})
		return
	}

	book, err := h.bookRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Book not found",
			Message: "The requested book could not be found",
		})
		return
	}

	response := dto.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Genre:       book.Genre,
		Purpose:     book.Purpose,
		Description: book.Description,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body dto.UpdateBookRequest true "Updated book information"
// @Success 200 {object} dto.BookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books/{id} [patch]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
		})
		return
	}

	var req dto.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Author != nil {
		updates["author"] = *req.Author
	}
	if req.Genre != nil {
		updates["genre"] = *req.Genre
	}
	if req.Purpose != nil {
		updates["purpose"] = *req.Purpose
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	book, err := h.bookRepo.Update(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Book not found",
			Message: "The requested book could not be found",
		})
		return
	}

	response := dto.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Genre:       book.Genre,
		Purpose:     book.Purpose,
		Description: book.Description,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} dto.BookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
		})
		return
	}

	book, err := h.bookRepo.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Book not found",
			Message: "The requested book could not be found",
		})
		return
	}

	response := dto.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Genre:       book.Genre,
		Purpose:     book.Purpose,
		Description: book.Description,
	}

	c.JSON(http.StatusOK, response)
}

// RecommendBook godoc
// @Summary Recommend a book
// @Description Get a book recommendation based on genre and purpose
// @Tags books
// @Accept json
// @Produce json
// @Param recommendation body dto.RecommendBookRequest true "Recommendation criteria"
// @Success 200 {object} dto.BookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books/recommend [post]
func (h *BookHandler) RecommendBook(c *gin.Context) {
	var req dto.RecommendBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	book, err := h.bookRepo.FindByGenreAndPurpose(req.Genre, req.Purpose)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "No recommendation found",
			Message: "No book found matching the criteria",
		})
		return
	}

	response := dto.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Genre:       book.Genre,
		Purpose:     book.Purpose,
		Description: book.Description,
	}

	c.JSON(http.StatusOK, response)
}

// AppError represents a custom error type
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
} 