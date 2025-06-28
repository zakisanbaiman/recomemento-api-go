package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"recomemento-api-go/database"
	"recomemento-api-go/dto"
	"recomemento-api-go/handlers"
	"recomemento-api-go/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// IntegrationTestSuite は統合テストスイートを定義
type IntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
}

// SetupSuite はテストスイート開始時に実行される
func (suite *IntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// テスト用データベース初期化
	dbPath := ":memory:"
	db, err := database.InitDatabase(dbPath)
	if err != nil {
		suite.T().Fatal("Failed to connect to test database:", err)
	}
	suite.db = db

	// リポジトリとハンドラーの初期化
	bookRepo := models.NewBookRepository(db)
	bookHandler := handlers.NewBookHandler(bookRepo)

	// ルーター設定
	r := gin.New()
	
	// CORS設定
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Recomemento API is running"})
	})

	// API routes
	api := r.Group("/")
	{
		api.POST("/books", bookHandler.CreateBook)
		api.GET("/books", bookHandler.GetAllBooks)
		api.GET("/books/:id", bookHandler.GetBookByID)
		api.PATCH("/books/:id", bookHandler.UpdateBook)
		api.DELETE("/books/:id", bookHandler.DeleteBook)
		api.POST("/books/recommend", bookHandler.RecommendBook)
	}

	suite.router = r
}

// SetupTest は各テスト前に実行される
func (suite *IntegrationTestSuite) SetupTest() {
	// 各テスト前にテーブルをクリア
	suite.db.Exec("DELETE FROM books")
}

// TestHealthCheck はヘルスチェックエンドポイントをテスト
func (suite *IntegrationTestSuite) TestHealthCheck() {
	// Act
	w := suite.performRequest("GET", "/health", nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "ok", response["status"])
	assert.Equal(suite.T(), "Recomemento API is running", response["message"])
}

// TestBookCRUD_FullFlow は本のCRUD操作の全フローをテスト
func (suite *IntegrationTestSuite) TestBookCRUD_FullFlow() {
	// 1. Create - 本を作成
	createReq := dto.CreateBookRequest{
		Title:       "Integration Test Book",
		Author:      "Test Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "This is a test book for integration testing",
	}

	body, _ := json.Marshal(createReq)
	w1 := suite.performRequest("POST", "/books", bytes.NewBuffer(body))
	
	assert.Equal(suite.T(), http.StatusCreated, w1.Code)
	
	var createdBook dto.BookResponse
	err := json.Unmarshal(w1.Body.Bytes(), &createdBook)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), createdBook.ID)
	assert.Equal(suite.T(), createReq.Title, createdBook.Title)
	assert.Equal(suite.T(), createReq.Author, createdBook.Author)
	assert.Equal(suite.T(), createReq.Genre, createdBook.Genre)

	bookID := createdBook.ID

	// 2. Read - 作成した本を取得
	w2 := suite.performRequest("GET", fmt.Sprintf("/books/%d", bookID), nil)
	
	assert.Equal(suite.T(), http.StatusOK, w2.Code)
	
	var retrievedBook dto.BookResponse
	err = json.Unmarshal(w2.Body.Bytes(), &retrievedBook)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdBook.ID, retrievedBook.ID)
	assert.Equal(suite.T(), createdBook.Title, retrievedBook.Title)
	assert.Equal(suite.T(), createdBook.Author, retrievedBook.Author)

	// 3. Update - 本を更新
	updateTitle := "Updated Integration Test Book"
	updateAuthor := "Updated Author"
	updateReq := dto.UpdateBookRequest{
		Title:  &updateTitle,
		Author: &updateAuthor,
	}

	body, _ = json.Marshal(updateReq)
	w3 := suite.performRequest("PATCH", fmt.Sprintf("/books/%d", bookID), bytes.NewBuffer(body))
	
	assert.Equal(suite.T(), http.StatusOK, w3.Code)
	
	var updatedBook dto.BookResponse
	err = json.Unmarshal(w3.Body.Bytes(), &updatedBook)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updateTitle, updatedBook.Title)
	assert.Equal(suite.T(), updateAuthor, updatedBook.Author)
	assert.Equal(suite.T(), createdBook.Genre, updatedBook.Genre) // 他のフィールドは変更されていない

	// 4. List - 全ての本を取得
	w4 := suite.performRequest("GET", "/books", nil)
	
	assert.Equal(suite.T(), http.StatusOK, w4.Code)
	
	var books []dto.BookResponse
	err = json.Unmarshal(w4.Body.Bytes(), &books)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), books, 1)
	assert.Equal(suite.T(), updateTitle, books[0].Title)
	assert.Equal(suite.T(), updateAuthor, books[0].Author)

	// 5. Delete - 本を削除
	w5 := suite.performRequest("DELETE", fmt.Sprintf("/books/%d", bookID), nil)
	
	assert.Equal(suite.T(), http.StatusOK, w5.Code)
	
	var deletedBook dto.BookResponse
	err = json.Unmarshal(w5.Body.Bytes(), &deletedBook)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updateTitle, deletedBook.Title)

	// 6. Verify deletion - 削除されたことを確認
	w6 := suite.performRequest("GET", fmt.Sprintf("/books/%d", bookID), nil)
	
	assert.Equal(suite.T(), http.StatusNotFound, w6.Code)
	
	// 7. Empty list - リストが空になったことを確認
	w7 := suite.performRequest("GET", "/books", nil)
	assert.Equal(suite.T(), http.StatusOK, w7.Code)
	
	var emptyBooks []dto.BookResponse
	err = json.Unmarshal(w7.Body.Bytes(), &emptyBooks)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), emptyBooks, 0)
}

// TestRecommendationFlow は推薦機能のフローをテスト
func (suite *IntegrationTestSuite) TestRecommendationFlow() {
	// 1. まず推薦用のデータを作成
	testBooks := []dto.CreateBookRequest{
		{
			Title: "Fiction Book 1", Author: "Author 1", Genre: "Fiction", 
			Purpose: "Entertainment", Description: "Fiction for entertainment",
		},
		{
			Title: "Tech Book 1", Author: "Author 2", Genre: "Technology", 
			Purpose: "Learning", Description: "Technology for learning",
		},
		{
			Title: "Fiction Book 2", Author: "Author 3", Genre: "Fiction", 
			Purpose: "Entertainment", Description: "Another fiction for entertainment",
		},
		{
			Title: "Business Book 1", Author: "Author 4", Genre: "Business", 
			Purpose: "Learning", Description: "Business for learning",
		},
	}

	for _, book := range testBooks {
		body, _ := json.Marshal(book)
		w := suite.performRequest("POST", "/books", bytes.NewBuffer(body))
		assert.Equal(suite.T(), http.StatusCreated, w.Code)
	}

	// 2. Fiction + Entertainment の推薦をリクエスト
	recommendReq := dto.RecommendBookRequest{
		Genre:   "Fiction",
		Purpose: "Entertainment",
	}

	body, _ := json.Marshal(recommendReq)
	w := suite.performRequest("POST", "/books/recommend", bytes.NewBuffer(body))

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var recommendedBook dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &recommendedBook)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Fiction", recommendedBook.Genre)
	assert.Equal(suite.T(), "Entertainment", recommendedBook.Purpose)
	assert.True(suite.T(), 
		recommendedBook.Title == "Fiction Book 1" || 
		recommendedBook.Title == "Fiction Book 2")

	// 3. Technology + Learning の推薦をリクエスト
	recommendReq2 := dto.RecommendBookRequest{
		Genre:   "Technology",
		Purpose: "Learning",
	}

	body, _ = json.Marshal(recommendReq2)
	w = suite.performRequest("POST", "/books/recommend", bytes.NewBuffer(body))

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	err = json.Unmarshal(w.Body.Bytes(), &recommendedBook)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Technology", recommendedBook.Genre)
	assert.Equal(suite.T(), "Learning", recommendedBook.Purpose)
	assert.Equal(suite.T(), "Tech Book 1", recommendedBook.Title)
}

// TestMultipleUsers は複数ユーザーの同時操作をシミュレート
func (suite *IntegrationTestSuite) TestMultipleUsers() {
	// User 1: Fiction の本を作成
	user1Book := dto.CreateBookRequest{
		Title: "User 1 Fiction", Author: "Author 1", Genre: "Fiction", 
		Purpose: "Entertainment", Description: "Fiction by user 1",
	}
	
	body, _ := json.Marshal(user1Book)
	w1 := suite.performRequest("POST", "/books", bytes.NewBuffer(body))
	assert.Equal(suite.T(), http.StatusCreated, w1.Code)
	
	var createdBook1 dto.BookResponse
	json.Unmarshal(w1.Body.Bytes(), &createdBook1)

	// User 2: Technology の本を作成
	user2Book := dto.CreateBookRequest{
		Title: "User 2 Tech", Author: "Author 2", Genre: "Technology", 
		Purpose: "Learning", Description: "Technology by user 2",
	}
	
	body, _ = json.Marshal(user2Book)
	w2 := suite.performRequest("POST", "/books", bytes.NewBuffer(body))
	assert.Equal(suite.T(), http.StatusCreated, w2.Code)
	
	var createdBook2 dto.BookResponse
	json.Unmarshal(w2.Body.Bytes(), &createdBook2)

	// User 1: Fiction の推薦を取得
	recommendReq1 := dto.RecommendBookRequest{
		Genre: "Fiction", Purpose: "Entertainment",
	}
	body, _ = json.Marshal(recommendReq1)
	w := suite.performRequest("POST", "/books/recommend", bytes.NewBuffer(body))
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// User 2: Technology の推薦を取得
	recommendReq2 := dto.RecommendBookRequest{
		Genre: "Technology", Purpose: "Learning",
	}
	body, _ = json.Marshal(recommendReq2)
	w = suite.performRequest("POST", "/books/recommend", bytes.NewBuffer(body))
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// All users: 全ての本を取得
	w = suite.performRequest("GET", "/books", nil)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var allBooks []dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &allBooks)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), allBooks, 2)
}

// TestErrorHandling はエラーハンドリングをテスト
func (suite *IntegrationTestSuite) TestErrorHandling() {
	// 1. 無効なJSONでの作成リクエスト
	w1 := suite.performRequest("POST", "/books", bytes.NewBufferString("invalid json"))
	assert.Equal(suite.T(), http.StatusBadRequest, w1.Code)

	var errorResp1 dto.ErrorResponse
	err := json.Unmarshal(w1.Body.Bytes(), &errorResp1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid request", errorResp1.Error)

	// 2. 存在しない本の取得
	w2 := suite.performRequest("GET", "/books/999", nil)
	assert.Equal(suite.T(), http.StatusNotFound, w2.Code)

	var errorResp2 dto.ErrorResponse
	err = json.Unmarshal(w2.Body.Bytes(), &errorResp2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Book not found", errorResp2.Error)

	// 3. 無効なIDでの取得
	w3 := suite.performRequest("GET", "/books/invalid-id", nil)
	assert.Equal(suite.T(), http.StatusBadRequest, w3.Code)

	var errorResp3 dto.ErrorResponse
	err = json.Unmarshal(w3.Body.Bytes(), &errorResp3)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid ID", errorResp3.Error)

	// 4. 存在しないジャンルでの推薦
	recommendReq := dto.RecommendBookRequest{
		Genre:   "NonExistentGenre",
		Purpose: "NonExistentPurpose",
	}
	body, _ := json.Marshal(recommendReq)
	w4 := suite.performRequest("POST", "/books/recommend", bytes.NewBuffer(body))
	assert.Equal(suite.T(), http.StatusNotFound, w4.Code)

	// 5. 必須フィールドが欠けた作成リクエスト
	incompleteReq := dto.CreateBookRequest{
		Title: "Only Title",
		// Author, Genre, Purpose, Description が欠けている
	}
	body, _ = json.Marshal(incompleteReq)
	w5 := suite.performRequest("POST", "/books", bytes.NewBuffer(body))
	assert.Equal(suite.T(), http.StatusBadRequest, w5.Code)

	// 6. 存在しない本の更新
	updateReq := dto.UpdateBookRequest{Title: stringPtr("Updated")}
	body, _ = json.Marshal(updateReq)
	w6 := suite.performRequest("PATCH", "/books/999", bytes.NewBuffer(body))
	assert.Equal(suite.T(), http.StatusNotFound, w6.Code)

	// 7. 存在しない本の削除
	w7 := suite.performRequest("DELETE", "/books/999", nil)
	assert.Equal(suite.T(), http.StatusNotFound, w7.Code)
}

// TestCORSHeaders はCORSヘッダーをテスト
func (suite *IntegrationTestSuite) TestCORSHeaders() {
	// 1. Preflight request (OPTIONS)
	req := httptest.NewRequest("OPTIONS", "/books", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(suite.T(), w.Header().Get("Access-Control-Allow-Methods"), "POST")

	// 2. Actual request with CORS headers
	createReq := dto.CreateBookRequest{
		Title: "CORS Test", Author: "Author", Genre: "Fiction", 
		Purpose: "Entertainment", Description: "CORS test book",
	}
	body, _ := json.Marshal(createReq)
	w = suite.performRequest("POST", "/books", bytes.NewBuffer(body))

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestDataPersistence はデータの永続性をテスト
func (suite *IntegrationTestSuite) TestDataPersistence() {
	// 1. 複数の本を作成
	books := []dto.CreateBookRequest{
		{Title: "Book 1", Author: "Author 1", Genre: "Fiction", Purpose: "Entertainment", Description: "Desc 1"},
		{Title: "Book 2", Author: "Author 2", Genre: "Technology", Purpose: "Learning", Description: "Desc 2"},
		{Title: "Book 3", Author: "Author 3", Genre: "Business", Purpose: "Learning", Description: "Desc 3"},
	}

	for _, book := range books {
		body, _ := json.Marshal(book)
		w := suite.performRequest("POST", "/books", bytes.NewBuffer(body))
		assert.Equal(suite.T(), http.StatusCreated, w.Code)
	}

	// 2. 全ての本を取得
	w := suite.performRequest("GET", "/books", nil)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var allBooks []dto.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &allBooks)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), allBooks, 3)

	// 3. 各本が正しく保存されているか確認
	titles := make([]string, len(allBooks))
	for i, book := range allBooks {
		titles[i] = book.Title
		assert.NotZero(suite.T(), book.ID)
		assert.NotEmpty(suite.T(), book.Author)
		assert.NotEmpty(suite.T(), book.Genre)
		assert.NotEmpty(suite.T(), book.Purpose)
	}
	
	assert.Contains(suite.T(), titles, "Book 1")
	assert.Contains(suite.T(), titles, "Book 2")
	assert.Contains(suite.T(), titles, "Book 3")
}

// ========== Helper Functions ==========

func (suite *IntegrationTestSuite) performRequest(method, url string, body *bytes.Buffer) *httptest.ResponseRecorder {
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

func stringPtr(s string) *string {
	return &s
}

// TestIntegrationTestSuite は統合テストスイートを実行
func TestIntegrationTestSuite(t *testing.T) {
	// 環境変数でテスト実行を制御
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.Skip("Integration tests skipped. Set RUN_INTEGRATION_TESTS=1 to run.")
	}
	
	suite.Run(t, new(IntegrationTestSuite))
} 