package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"recomemento-api-go/dto"
	"recomemento-api-go/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// BookHandlerExtendedTestSuite はテストスイートを定義
type BookHandlerExtendedTestSuite struct {
	suite.Suite
	handler  *BookHandler
	mockRepo *MockExtendedBookDatabase
	router   *gin.Engine
}

// MockExtendedBookDatabase は拡張されたモックデータベース
type MockExtendedBookDatabase struct {
	mock.Mock
}

func (m *MockExtendedBookDatabase) Create(book *models.Book) error {
	args := m.Called(book)
	if args.Error(0) == nil {
		book.ID = 1 // Simulate auto-generated ID
	}
	return args.Error(0)
}

func (m *MockExtendedBookDatabase) GetAll() ([]models.Book, error) {
	args := m.Called()
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockExtendedBookDatabase) GetByID(id uint) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockExtendedBookDatabase) Update(id uint, updates map[string]interface{}) (*models.Book, error) {
	args := m.Called(id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockExtendedBookDatabase) Delete(id uint) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockExtendedBookDatabase) FindByGenreAndPurpose(genre, purpose string) (*models.Book, error) {
	args := m.Called(genre, purpose)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

// SetupTest は各テスト前に実行される
func (suite *BookHandlerExtendedTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockRepo = new(MockExtendedBookDatabase)
	suite.handler = NewBookHandler(suite.mockRepo)
	suite.router = gin.New()
	
	// ルート設定
	suite.router.POST("/books", suite.handler.CreateBook)
	suite.router.GET("/books", suite.handler.GetAllBooks)
	suite.router.GET("/books/:id", suite.handler.GetBookByID)
	suite.router.PATCH("/books/:id", suite.handler.UpdateBook)
	suite.router.DELETE("/books/:id", suite.handler.DeleteBook)
	suite.router.POST("/books/recommend", suite.handler.RecommendBook)
}

// ========== CreateBook Tests ==========

func (suite *BookHandlerExtendedTestSuite) TestCreateBook_Success() {
	// Arrange
	req := dto.CreateBookRequest{
		Title:       "Test Book",
		Author:      "Test Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Test Description",
	}
	
	suite.mockRepo.On("Create", mock.AnythingOfType("*models.Book")).Return(nil)

	// Act
	body, _ := json.Marshal(req)
	w := suite.performRequest("POST", "/books", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	
	var response dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), req.Title, response.Title)
	assert.Equal(suite.T(), req.Author, response.Author)
	assert.Equal(suite.T(), uint(1), response.ID)
	
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *BookHandlerExtendedTestSuite) TestCreateBook_ValidationError_MissingTitle() {
	// Arrange - タイトルが欠けているリクエスト
	req := dto.CreateBookRequest{
		// Title: "", // 欠けている
		Author:      "Test Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Test Description",
	}

	// Act
	body, _ := json.Marshal(req)
	w := suite.performRequest("POST", "/books", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid request", response.Error)
	assert.Contains(suite.T(), response.Message, "Title")
}

func (suite *BookHandlerExtendedTestSuite) TestCreateBook_ValidationError_AllFieldsMissing() {
	// Arrange - 全てのフィールドが欠けているリクエスト
	req := dto.CreateBookRequest{}

	// Act
	body, _ := json.Marshal(req)
	w := suite.performRequest("POST", "/books", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid request", response.Error)
}

func (suite *BookHandlerExtendedTestSuite) TestCreateBook_InvalidJSON() {
	// Act - 無効なJSONを送信
	w := suite.performRequest("POST", "/books", bytes.NewBufferString("invalid json"))

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid request", response.Error)
}

func (suite *BookHandlerExtendedTestSuite) TestCreateBook_DatabaseError() {
	// Arrange
	req := dto.CreateBookRequest{
		Title:       "Test Book",
		Author:      "Test Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Test Description",
	}
	
	suite.mockRepo.On("Create", mock.AnythingOfType("*models.Book")).Return(errors.New("database error"))

	// Act
	body, _ := json.Marshal(req)
	w := suite.performRequest("POST", "/books", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Failed to create book", response.Error)
	assert.Equal(suite.T(), "database error", response.Message)
}

// ========== GetBookByID Tests ==========

func (suite *BookHandlerExtendedTestSuite) TestGetBookByID_Success() {
	// Arrange
	expectedBook := &models.Book{
		ID: 1, Title: "Test Book", Author: "Test Author",
		Genre: "Fiction", Purpose: "Entertainment", Description: "Test Description",
	}
	
	suite.mockRepo.On("GetByID", uint(1)).Return(expectedBook, nil)

	// Act
	w := suite.performRequest("GET", "/books/1", nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedBook.Title, response.Title)
	assert.Equal(suite.T(), expectedBook.ID, response.ID)
}

func (suite *BookHandlerExtendedTestSuite) TestGetBookByID_InvalidID() {
	// Act
	w := suite.performRequest("GET", "/books/invalid", nil)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid ID", response.Error)
	assert.Equal(suite.T(), "ID must be a valid number", response.Message)
}

func (suite *BookHandlerExtendedTestSuite) TestGetBookByID_NotFound() {
	// Arrange
	suite.mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("record not found"))

	// Act
	w := suite.performRequest("GET", "/books/999", nil)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Book not found", response.Error)
}

func (suite *BookHandlerExtendedTestSuite) TestGetBookByID_NegativeID() {
	// Act
	w := suite.performRequest("GET", "/books/-1", nil)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// ========== GetAllBooks Tests ==========

func (suite *BookHandlerExtendedTestSuite) TestGetAllBooks_Success() {
	// Arrange
	expectedBooks := []models.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 1"},
		{ID: 2, Title: "Book 2", Author: "Author 2", Genre: "Technology", Purpose: "Learning", Description: "Description 2"},
	}
	
	suite.mockRepo.On("GetAll").Return(expectedBooks, nil)

	// Act
	w := suite.performRequest("GET", "/books", nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response []dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response, 2)
	assert.Equal(suite.T(), "Book 1", response[0].Title)
	assert.Equal(suite.T(), "Book 2", response[1].Title)
}

func (suite *BookHandlerExtendedTestSuite) TestGetAllBooks_EmptyResult() {
	// Arrange
	emptyBooks := []models.Book{}
	suite.mockRepo.On("GetAll").Return(emptyBooks, nil)

	// Act
	w := suite.performRequest("GET", "/books", nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response []dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response, 0)
}

func (suite *BookHandlerExtendedTestSuite) TestGetAllBooks_DatabaseError() {
	// Arrange
	suite.mockRepo.On("GetAll").Return([]models.Book{}, errors.New("database connection failed"))

	// Act
	w := suite.performRequest("GET", "/books", nil)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Failed to get books", response.Error)
}

// ========== UpdateBook Tests ==========

func (suite *BookHandlerExtendedTestSuite) TestUpdateBook_Success() {
	// Arrange
	updatedBook := &models.Book{
		ID: 1, Title: "Updated Title", Author: "Original Author",
		Genre: "Fiction", Purpose: "Entertainment", Description: "Original Description",
	}
	
	updates := map[string]interface{}{"title": "Updated Title"}
	suite.mockRepo.On("Update", uint(1), updates).Return(updatedBook, nil)

	updateReq := dto.UpdateBookRequest{
		Title: stringPointer("Updated Title"),
	}

	// Act
	body, _ := json.Marshal(updateReq)
	w := suite.performRequest("PATCH", "/books/1", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Title", response.Title)
	assert.Equal(suite.T(), "Original Author", response.Author)
}

func (suite *BookHandlerExtendedTestSuite) TestUpdateBook_NotFound() {
	// Arrange
	suite.mockRepo.On("Update", uint(999), mock.Anything).Return(nil, errors.New("record not found"))

	updateReq := dto.UpdateBookRequest{
		Title: stringPointer("Updated Title"),
	}

	// Act
	body, _ := json.Marshal(updateReq)
	w := suite.performRequest("PATCH", "/books/999", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

// ========== DeleteBook Tests ==========

func (suite *BookHandlerExtendedTestSuite) TestDeleteBook_Success() {
	// Arrange
	deletedBook := &models.Book{
		ID: 1, Title: "Book to Delete", Author: "Author",
		Genre: "Fiction", Purpose: "Entertainment", Description: "Description",
	}
	
	suite.mockRepo.On("Delete", uint(1)).Return(deletedBook, nil)

	// Act
	w := suite.performRequest("DELETE", "/books/1", nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), deletedBook.Title, response.Title)
}

// ========== RecommendBook Tests ==========

func (suite *BookHandlerExtendedTestSuite) TestRecommendBook_Success() {
	// Arrange
	recommendedBook := &models.Book{
		ID: 1, Title: "Recommended Book", Author: "Author",
		Genre: "Fiction", Purpose: "Entertainment", Description: "Description",
	}
	
	suite.mockRepo.On("FindByGenreAndPurpose", "Fiction", "Entertainment").Return(recommendedBook, nil)

	recommendReq := dto.RecommendBookRequest{
		Genre:   "Fiction",
		Purpose: "Entertainment",
	}

	// Act
	body, _ := json.Marshal(recommendReq)
	w := suite.performRequest("POST", "/books/recommend", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Recommended Book", response.Title)
}

func (suite *BookHandlerExtendedTestSuite) TestRecommendBook_NotFound() {
	// Arrange
	suite.mockRepo.On("FindByGenreAndPurpose", "NonExistent", "Purpose").Return(nil, errors.New("record not found"))

	recommendReq := dto.RecommendBookRequest{
		Genre:   "NonExistent",
		Purpose: "Purpose",
	}

	// Act
	body, _ := json.Marshal(recommendReq)
	w := suite.performRequest("POST", "/books/recommend", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *BookHandlerExtendedTestSuite) TestRecommendBook_ValidationError() {
	// Arrange - 必須フィールドが欠けている
	recommendReq := dto.RecommendBookRequest{
		Genre: "Fiction",
		// Purpose が欠けている
	}

	// Act
	body, _ := json.Marshal(recommendReq)
	w := suite.performRequest("POST", "/books/recommend", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// ========== Helper Functions ==========

func (suite *BookHandlerExtendedTestSuite) performRequest(method, url string, body *bytes.Buffer) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, url, body)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func stringPointer(s string) *string {
	return &s
}

// TestSuite実行
func TestBookHandlerExtendedTestSuite(t *testing.T) {
	suite.Run(t, new(BookHandlerExtendedTestSuite))
} 