package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"recomemento-api-go/dto"
	"recomemento-api-go/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBookDatabase is a mock implementation of BookDatabase interface
type MockBookDatabase struct {
	mock.Mock
}

func (m *MockBookDatabase) Create(book *models.Book) error {
	args := m.Called(book)
	book.ID = 1 // Simulate auto-generated ID
	return args.Error(0)
}

func (m *MockBookDatabase) GetAll() ([]models.Book, error) {
	args := m.Called()
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBookDatabase) GetByID(id uint) (*models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookDatabase) Update(id uint, updates map[string]interface{}) (*models.Book, error) {
	args := m.Called(id, updates)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookDatabase) Delete(id uint) (*models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookDatabase) FindByGenreAndPurpose(genre, purpose string) (*models.Book, error) {
	args := m.Called(genre, purpose)
	return args.Get(0).(*models.Book), args.Error(1)
}

func TestCreateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockBookDatabase)
	handler := NewBookHandler(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*models.Book")).Return(nil)

	createRequest := dto.CreateBookRequest{
		Title:       "Test Book",
		Author:      "Test Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Test Description",
	}

	jsonValue, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateBook(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Book", response.Title)
	assert.Equal(t, "Test Author", response.Author)

	mockRepo.AssertExpectations(t)
}

func TestGetAllBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockBookDatabase)
	handler := NewBookHandler(mockRepo)

	expectedBooks := []models.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 1"},
		{ID: 2, Title: "Book 2", Author: "Author 2", Genre: "Technology", Purpose: "Learning", Description: "Description 2"},
	}

	mockRepo.On("GetAll").Return(expectedBooks, nil)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetAllBooks(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Book 1", response[0].Title)
	assert.Equal(t, "Book 2", response[1].Title)

	mockRepo.AssertExpectations(t)
}

func TestRecommendBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockBookDatabase)
	handler := NewBookHandler(mockRepo)

	expectedBook := &models.Book{
		ID: 1, Title: "Recommended Book", Author: "Author", 
		Genre: "Fiction", Purpose: "Entertainment", Description: "Description",
	}

	mockRepo.On("FindByGenreAndPurpose", "Fiction", "Entertainment").Return(expectedBook, nil)

	recommendRequest := dto.RecommendBookRequest{
		Genre:   "Fiction",
		Purpose: "Entertainment",
	}

	jsonValue, _ := json.Marshal(recommendRequest)
	req, _ := http.NewRequest("POST", "/books/recommend", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.RecommendBook(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Recommended Book", response.Title)

	mockRepo.AssertExpectations(t)
} 