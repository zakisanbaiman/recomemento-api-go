package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BookRepositoryTestSuite はリポジトリテストスイートを定義
type BookRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo BookDatabase
}

// SetupSuite はテストスイート開始時に実行される
func (suite *BookRepositoryTestSuite) SetupSuite() {
	// テスト用のインメモリSQLiteデータベースを作成
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	
	if err != nil {
		suite.T().Fatal("Failed to connect to test database:", err)
	}

	// マイグレーション実行
	err = db.AutoMigrate(&Book{})
	if err != nil {
		suite.T().Fatal("Failed to migrate test database:", err)
	}

	suite.db = db
	suite.repo = NewBookRepository(db)
}

// SetupTest は各テスト前に実行される
func (suite *BookRepositoryTestSuite) SetupTest() {
	// 各テスト前にテーブルをクリア
	suite.db.Exec("DELETE FROM books")
}

// ========== Create Tests ==========

func (suite *BookRepositoryTestSuite) TestCreate_Success() {
	// Arrange
	book := &Book{
		Title:       "Test Book",
		Author:      "Test Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Test Description",
	}

	// Act
	err := suite.repo.Create(book)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), book.ID)

	// データベースに実際に保存されているか確認
	var savedBook Book
	suite.db.First(&savedBook, book.ID)
	assert.Equal(suite.T(), book.Title, savedBook.Title)
	assert.Equal(suite.T(), book.Author, savedBook.Author)
	assert.Equal(suite.T(), book.Genre, savedBook.Genre)
	assert.Equal(suite.T(), book.Purpose, savedBook.Purpose)
	assert.Equal(suite.T(), book.Description, savedBook.Description)
}

func (suite *BookRepositoryTestSuite) TestCreate_MultipleBooks() {
	// Arrange
	books := []*Book{
		{Title: "Book 1", Author: "Author 1", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 1"},
		{Title: "Book 2", Author: "Author 2", Genre: "Technology", Purpose: "Learning", Description: "Description 2"},
		{Title: "Book 3", Author: "Author 3", Genre: "Business", Purpose: "Learning", Description: "Description 3"},
	}

	// Act & Assert
	for _, book := range books {
		err := suite.repo.Create(book)
		assert.NoError(suite.T(), err)
		assert.NotZero(suite.T(), book.ID)
	}

	// 全ての本が保存されているか確認
	var count int64
	suite.db.Model(&Book{}).Count(&count)
	assert.Equal(suite.T(), int64(3), count)
}

func (suite *BookRepositoryTestSuite) TestCreate_WithUnicodeCharacters() {
	// Arrange - 日本語を含む本
	book := &Book{
		Title:       "吾輩は猫である",
		Author:      "夏目漱石",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "猫の視点から描かれた小説",
	}

	// Act
	err := suite.repo.Create(book)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), book.ID)

	// データベースから取得して確認
	var savedBook Book
	suite.db.First(&savedBook, book.ID)
	assert.Equal(suite.T(), "吾輩は猫である", savedBook.Title)
	assert.Equal(suite.T(), "夏目漱石", savedBook.Author)
}

// ========== GetByID Tests ==========

func (suite *BookRepositoryTestSuite) TestGetByID_Success() {
	// Arrange - テストデータを直接データベースに挿入
	book := Book{
		Title:       "Test Book",
		Author:      "Test Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Test Description",
	}
	suite.db.Create(&book)

	// Act
	result, err := suite.repo.GetByID(book.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), book.ID, result.ID)
	assert.Equal(suite.T(), book.Title, result.Title)
	assert.Equal(suite.T(), book.Author, result.Author)
	assert.Equal(suite.T(), book.Genre, result.Genre)
	assert.Equal(suite.T(), book.Purpose, result.Purpose)
	assert.Equal(suite.T(), book.Description, result.Description)
}

func (suite *BookRepositoryTestSuite) TestGetByID_NotFound() {
	// Act
	result, err := suite.repo.GetByID(999)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BookRepositoryTestSuite) TestGetByID_ZeroID() {
	// Act
	result, err := suite.repo.GetByID(0)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// ========== GetAll Tests ==========

func (suite *BookRepositoryTestSuite) TestGetAll_Success() {
	// Arrange - 複数のテストデータを挿入
	books := []Book{
		{Title: "Book 1", Author: "Author 1", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 1"},
		{Title: "Book 2", Author: "Author 2", Genre: "Technology", Purpose: "Learning", Description: "Description 2"},
		{Title: "Book 3", Author: "Author 3", Genre: "Business", Purpose: "Learning", Description: "Description 3"},
	}

	for _, book := range books {
		suite.db.Create(&book)
	}

	// Act
	result, err := suite.repo.GetAll()

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)
	
	// 順序は保証されないので、タイトルでソートして確認
	titles := make([]string, len(result))
	for i, book := range result {
		titles[i] = book.Title
	}
	assert.Contains(suite.T(), titles, "Book 1")
	assert.Contains(suite.T(), titles, "Book 2")
	assert.Contains(suite.T(), titles, "Book 3")
}

func (suite *BookRepositoryTestSuite) TestGetAll_Empty() {
	// Act - データが空の状態で取得
	result, err := suite.repo.GetAll()

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
}

func (suite *BookRepositoryTestSuite) TestGetAll_LargeDataset() {
	// Arrange - 大量のデータを挿入
	booksCount := 100
	for i := 1; i <= booksCount; i++ {
		book := Book{
			Title:       fmt.Sprintf("Book %d", i),
			Author:      fmt.Sprintf("Author %d", i),
			Genre:       "Fiction",
			Purpose:     "Entertainment",
			Description: fmt.Sprintf("Description %d", i),
		}
		suite.db.Create(&book)
	}

	// Act
	result, err := suite.repo.GetAll()

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, booksCount)
}

// ========== Update Tests ==========

func (suite *BookRepositoryTestSuite) TestUpdate_Success() {
	// Arrange
	book := Book{
		Title:       "Original Title",
		Author:      "Original Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Original Description",
	}
	suite.db.Create(&book)

	updates := map[string]interface{}{
		"title":  "Updated Title",
		"author": "Updated Author",
	}

	// Act
	result, err := suite.repo.Update(book.ID, updates)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Updated Title", result.Title)
	assert.Equal(suite.T(), "Updated Author", result.Author)
	assert.Equal(suite.T(), book.Genre, result.Genre) // 更新されていないフィールドは元のまま
	assert.Equal(suite.T(), book.Purpose, result.Purpose)
	assert.Equal(suite.T(), book.Description, result.Description)

	// データベースからも確認
	var updatedBookFromDB Book
	suite.db.First(&updatedBookFromDB, book.ID)
	assert.Equal(suite.T(), "Updated Title", updatedBookFromDB.Title)
	assert.Equal(suite.T(), "Updated Author", updatedBookFromDB.Author)
}

func (suite *BookRepositoryTestSuite) TestUpdate_PartialUpdate() {
	// Arrange
	book := Book{
		Title:       "Original Title",
		Author:      "Original Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Original Description",
	}
	suite.db.Create(&book)

	updates := map[string]interface{}{
		"title": "Updated Title Only",
	}

	// Act
	result, err := suite.repo.Update(book.ID, updates)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Title Only", result.Title)
	assert.Equal(suite.T(), "Original Author", result.Author) // 変更されていない
}

func (suite *BookRepositoryTestSuite) TestUpdate_NotFound() {
	// Arrange
	updates := map[string]interface{}{
		"title": "Updated Title",
	}

	// Act
	result, err := suite.repo.Update(999, updates)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BookRepositoryTestSuite) TestUpdate_EmptyUpdates() {
	// Arrange
	book := Book{
		Title:       "Original Title",
		Author:      "Original Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Original Description",
	}
	suite.db.Create(&book)

	updates := map[string]interface{}{}

	// Act
	result, err := suite.repo.Update(book.ID, updates)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// 何も変更されていないことを確認
	assert.Equal(suite.T(), book.Title, result.Title)
	assert.Equal(suite.T(), book.Author, result.Author)
}

// ========== Delete Tests ==========

func (suite *BookRepositoryTestSuite) TestDelete_Success() {
	// Arrange
	book := Book{
		Title:       "Book to Delete",
		Author:      "Author",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Description",
	}
	suite.db.Create(&book)

	// Act
	result, err := suite.repo.Delete(book.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), book.Title, result.Title)

	// データベースから実際に削除されているか確認
	var deletedBook Book
	err = suite.db.First(&deletedBook, book.ID).Error
	assert.Error(suite.T(), err) // レコードが見つからないエラーが発生するはず
}

func (suite *BookRepositoryTestSuite) TestDelete_NotFound() {
	// Act
	result, err := suite.repo.Delete(999)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BookRepositoryTestSuite) TestDelete_MultipleAndRetrieve() {
	// Arrange - 複数の本を作成
	books := []Book{
		{Title: "Book 1", Author: "Author 1", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 1"},
		{Title: "Book 2", Author: "Author 2", Genre: "Technology", Purpose: "Learning", Description: "Description 2"},
		{Title: "Book 3", Author: "Author 3", Genre: "Business", Purpose: "Learning", Description: "Description 3"},
	}

	for i := range books {
		suite.db.Create(&books[i])
	}

	// Act - 真ん中の本を削除
	result, err := suite.repo.Delete(books[1].ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Book 2", result.Title)

	// 残りの本が2冊あることを確認
	allBooks, err := suite.repo.GetAll()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), allBooks, 2)
}

// ========== FindByGenreAndPurpose Tests ==========

func (suite *BookRepositoryTestSuite) TestFindByGenreAndPurpose_Success() {
	// Arrange
	books := []Book{
		{Title: "Fiction Book", Author: "Author 1", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 1"},
		{Title: "Tech Book", Author: "Author 2", Genre: "Technology", Purpose: "Learning", Description: "Description 2"},
		{Title: "Another Fiction", Author: "Author 3", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 3"},
		{Title: "Business Book", Author: "Author 4", Genre: "Business", Purpose: "Learning", Description: "Description 4"},
	}

	for _, book := range books {
		suite.db.Create(&book)
	}

	// Act
	result, err := suite.repo.FindByGenreAndPurpose("Fiction", "Entertainment")

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Fiction", result.Genre)
	assert.Equal(suite.T(), "Entertainment", result.Purpose)
	// どちらの本が返されるかは不定だが、条件に合致していることを確認
	assert.True(suite.T(), result.Title == "Fiction Book" || result.Title == "Another Fiction")
}

func (suite *BookRepositoryTestSuite) TestFindByGenreAndPurpose_NotFound() {
	// Arrange
	book := Book{
		Title: "Fiction Book", Author: "Author 1", Genre: "Fiction", 
		Purpose: "Entertainment", Description: "Description 1",
	}
	suite.db.Create(&book)

	// Act
	result, err := suite.repo.FindByGenreAndPurpose("NonExistent", "Purpose")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BookRepositoryTestSuite) TestFindByGenreAndPurpose_CaseSensitive() {
	// Arrange
	book := Book{
		Title: "Fiction Book", Author: "Author 1", Genre: "Fiction", 
		Purpose: "Entertainment", Description: "Description 1",
	}
	suite.db.Create(&book)

	// Act - 大文字小文字が異なる場合
	result, err := suite.repo.FindByGenreAndPurpose("fiction", "entertainment")

	// Assert - SQLiteは大文字小文字を区別するので見つからない
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BookRepositoryTestSuite) TestFindByGenreAndPurpose_MultipleMatches() {
	// Arrange - 同じ条件の本を複数作成
	books := []Book{
		{Title: "Fiction Book 1", Author: "Author 1", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 1"},
		{Title: "Fiction Book 2", Author: "Author 2", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 2"},
		{Title: "Fiction Book 3", Author: "Author 3", Genre: "Fiction", Purpose: "Entertainment", Description: "Description 3"},
	}

	for _, book := range books {
		suite.db.Create(&book)
	}

	// Act
	result, err := suite.repo.FindByGenreAndPurpose("Fiction", "Entertainment")

	// Assert - 最初に見つかった1件が返される
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Fiction", result.Genre)
	assert.Equal(suite.T(), "Entertainment", result.Purpose)
	assert.Contains(suite.T(), result.Title, "Fiction Book")
}

// ========== Edge Cases ==========

func (suite *BookRepositoryTestSuite) TestRepository_WithSpecialCharacters() {
	// Arrange - 特殊文字を含む本
	book := &Book{
		Title:       `Test "Book" with 'Quotes' & <Tags>`,
		Author:      "Author with @#$%^&*() symbols",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "Description with\nnewlines\tand\ttabs",
	}

	// Act
	err := suite.repo.Create(book)
	assert.NoError(suite.T(), err)

	result, err := suite.repo.GetByID(book.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), book.Title, result.Title)
	assert.Equal(suite.T(), book.Author, result.Author)
	assert.Equal(suite.T(), book.Description, result.Description)
}

func (suite *BookRepositoryTestSuite) TestRepository_WithEmptyStrings() {
	// Arrange - 空文字列を含む本（バリデーションは上位層で行うため、ここでは許可）
	book := &Book{
		Title:       "",
		Author:      "",
		Genre:       "Fiction",
		Purpose:     "Entertainment",
		Description: "",
	}

	// Act
	err := suite.repo.Create(book)

	// Assert - データベースレベルでは空文字列も許可される
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), book.ID)
}

// TestBookRepositoryTestSuite はリポジトリテストスイートを実行
func TestBookRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BookRepositoryTestSuite))
} 