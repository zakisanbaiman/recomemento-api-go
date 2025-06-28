package testutil

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"recomemento-api-go/dto"
	"recomemento-api-go/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BookFactory はテスト用の本データを生成するファクトリー
type BookFactory struct {
	counter int
	rand    *rand.Rand
}

// NewBookFactory は新しいBookFactoryを作成
func NewBookFactory() *BookFactory {
	return &BookFactory{
		counter: 0,
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreateBook はテスト用の本データを作成
func (f *BookFactory) CreateBook(overrides ...func(*models.Book)) *models.Book {
	f.counter++
	
	book := &models.Book{
		Title:       fmt.Sprintf("Test Book %d", f.counter),
		Author:      fmt.Sprintf("Test Author %d", f.counter),
		Genre:       f.randomGenre(),
		Purpose:     f.randomPurpose(),
		Description: fmt.Sprintf("Test Description for book %d", f.counter),
	}

	// オーバーライドの適用
	for _, override := range overrides {
		override(book)
	}

	return book
}

// CreateBookRequest はテスト用のCreateBookRequestを作成
func (f *BookFactory) CreateBookRequest(overrides ...func(*dto.CreateBookRequest)) dto.CreateBookRequest {
	f.counter++
	
	req := dto.CreateBookRequest{
		Title:       fmt.Sprintf("Test Book %d", f.counter),
		Author:      fmt.Sprintf("Test Author %d", f.counter),
		Genre:       f.randomGenre(),
		Purpose:     f.randomPurpose(),
		Description: fmt.Sprintf("Test Description for book %d", f.counter),
	}

	// オーバーライドの適用
	for _, override := range overrides {
		override(&req)
	}

	return req
}

// CreateBooks は複数の本データを作成
func (f *BookFactory) CreateBooks(count int) []models.Book {
	books := make([]models.Book, count)
	for i := 0; i < count; i++ {
		book := f.CreateBook()
		books[i] = *book
	}
	return books
}

// CreateBookRequests は複数のCreateBookRequestを作成
func (f *BookFactory) CreateBookRequests(count int) []dto.CreateBookRequest {
	requests := make([]dto.CreateBookRequest, count)
	for i := 0; i < count; i++ {
		requests[i] = f.CreateBookRequest()
	}
	return requests
}

// CreateFictionBook はフィクション本を作成
func (f *BookFactory) CreateFictionBook() *models.Book {
	return f.CreateBook(WithGenre("Fiction"), WithPurpose("Entertainment"))
}

// CreateTechBook は技術本を作成
func (f *BookFactory) CreateTechBook() *models.Book {
	return f.CreateBook(WithGenre("Technology"), WithPurpose("Learning"))
}

// CreateBusinessBook はビジネス本を作成
func (f *BookFactory) CreateBusinessBook() *models.Book {
	return f.CreateBook(WithGenre("Business"), WithPurpose("Learning"))
}

// ========== Override Functions ==========

// WithGenre はジャンルを設定するオーバーライド関数
func WithGenre(genre string) func(*models.Book) {
	return func(book *models.Book) {
		book.Genre = genre
	}
}

// WithPurpose は目的を設定するオーバーライド関数
func WithPurpose(purpose string) func(*models.Book) {
	return func(book *models.Book) {
		book.Purpose = purpose
	}
}

// WithTitle はタイトルを設定するオーバーライド関数
func WithTitle(title string) func(*models.Book) {
	return func(book *models.Book) {
		book.Title = title
	}
}

// WithAuthor は著者を設定するオーバーライド関数
func WithAuthor(author string) func(*models.Book) {
	return func(book *models.Book) {
		book.Author = author
	}
}

// WithDescription は説明を設定するオーバーライド関数
func WithDescription(description string) func(*models.Book) {
	return func(book *models.Book) {
		book.Description = description
	}
}

// WithID はIDを設定するオーバーライド関数
func WithID(id uint) func(*models.Book) {
	return func(book *models.Book) {
		book.ID = id
	}
}

// ========== CreateBookRequest Override Functions ==========

// WithRequestGenre はCreateBookRequestのジャンルを設定
func WithRequestGenre(genre string) func(*dto.CreateBookRequest) {
	return func(req *dto.CreateBookRequest) {
		req.Genre = genre
	}
}

// WithRequestPurpose はCreateBookRequestの目的を設定
func WithRequestPurpose(purpose string) func(*dto.CreateBookRequest) {
	return func(req *dto.CreateBookRequest) {
		req.Purpose = purpose
	}
}

// WithRequestTitle はCreateBookRequestのタイトルを設定
func WithRequestTitle(title string) func(*dto.CreateBookRequest) {
	return func(req *dto.CreateBookRequest) {
		req.Title = title
	}
}

// WithRequestAuthor はCreateBookRequestの著者を設定
func WithRequestAuthor(author string) func(*dto.CreateBookRequest) {
	return func(req *dto.CreateBookRequest) {
		req.Author = author
	}
}

// ========== Test Database Helper ==========

// TestDatabase はテスト用データベースを提供
type TestDatabase struct {
	DB *gorm.DB
}

// NewTestDatabase は新しいテスト用データベースを作成
func NewTestDatabase(t *testing.T) *TestDatabase {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}

	// マイグレーション実行
	err = db.AutoMigrate(&models.Book{})
	if err != nil {
		t.Fatal("Failed to migrate test database:", err)
	}

	return &TestDatabase{DB: db}
}

// CleanUp はテストデータをクリーンアップ
func (td *TestDatabase) CleanUp() {
	td.DB.Exec("DELETE FROM books")
}

// SeedBooks は複数の本データをデータベースに挿入
func (td *TestDatabase) SeedBooks(books []models.Book) error {
	for _, book := range books {
		if err := td.DB.Create(&book).Error; err != nil {
			return err
		}
	}
	return nil
}

// SeedBook は単一の本をデータベースに挿入
func (td *TestDatabase) SeedBook(book *models.Book) error {
	return td.DB.Create(book).Error
}

// CountBooks はデータベース内の本の数を取得
func (td *TestDatabase) CountBooks() int64 {
	var count int64
	td.DB.Model(&models.Book{}).Count(&count)
	return count
}

// FindBookByTitle はタイトルで本を検索
func (td *TestDatabase) FindBookByTitle(title string) (*models.Book, error) {
	var book models.Book
	err := td.DB.Where("title = ?", title).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// ========== Assertion Helpers ==========

// AssertJSONEqual はJSONレスポンスが期待値と一致するかを検証
func AssertJSONEqual(t *testing.T, expected, actual interface{}) {
	expectedJSON, err := json.Marshal(expected)
	assert.NoError(t, err)

	actualJSON, err := json.Marshal(actual)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}

// AssertBookEqual は本オブジェクトが等しいかを検証（IDを除く）
func AssertBookEqual(t *testing.T, expected, actual *models.Book) {
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Author, actual.Author)
	assert.Equal(t, expected.Genre, actual.Genre)
	assert.Equal(t, expected.Purpose, actual.Purpose)
	assert.Equal(t, expected.Description, actual.Description)
}

// AssertBookResponseEqual はBookResponseが等しいかを検証
func AssertBookResponseEqual(t *testing.T, expected *models.Book, actual *dto.BookResponse) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Author, actual.Author)
	assert.Equal(t, expected.Genre, actual.Genre)
	assert.Equal(t, expected.Purpose, actual.Purpose)
	assert.Equal(t, expected.Description, actual.Description)
}

// AssertErrorResponse はエラーレスポンスを検証
func AssertErrorResponse(t *testing.T, body []byte, expectedError, expectedMessage string) {
	var errorResp dto.ErrorResponse
	err := json.Unmarshal(body, &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedError, errorResp.Error)
	if expectedMessage != "" {
		assert.Equal(t, expectedMessage, errorResp.Message)
	}
}

// ========== Random Data Generators ==========

// RandomString はランダムな文字列を生成
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomInt は指定された範囲内のランダムな整数を生成
func RandomInt(min, max int) int {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return seededRand.Intn(max-min+1) + min
}

// RandomBool はランダムなブール値を生成
func RandomBool() bool {
	return RandomInt(0, 1) == 1
}

// ========== Predefined Data ==========

// GenreOptions はテスト用のジャンルオプション
var GenreOptions = []string{
	"Fiction", "Non-Fiction", "Technology", "Business", "Science", 
	"History", "Biography", "Self-Help", "Travel", "Cooking",
}

// PurposeOptions はテスト用の目的オプション
var PurposeOptions = []string{
	"Entertainment", "Learning", "Research", "Reference", "Inspiration",
}

// SampleAuthors はテスト用の著者リスト
var SampleAuthors = []string{
	"山田太郎", "田中花子", "佐藤次郎", "鈴木美香", "高橋一郎",
	"John Smith", "Jane Doe", "Robert Johnson", "Emily Davis", "Michael Brown",
}

// SampleTitles はテスト用のタイトルリスト
var SampleTitles = []string{
	"素晴らしい本", "技術の未来", "ビジネス戦略", "人生の教訓", "冒険の物語",
	"The Great Adventure", "Tech Revolution", "Business Mastery", "Life Lessons", "Creative Writing",
}

// ========== Factory Helper Methods ==========

func (f *BookFactory) randomGenre() string {
	return GenreOptions[f.rand.Intn(len(GenreOptions))]
}

func (f *BookFactory) randomPurpose() string {
	return PurposeOptions[f.rand.Intn(len(PurposeOptions))]
}

func (f *BookFactory) randomAuthor() string {
	return SampleAuthors[f.rand.Intn(len(SampleAuthors))]
}

func (f *BookFactory) randomTitle() string {
	return SampleTitles[f.rand.Intn(len(SampleTitles))]
}

// RandomGenre はランダムなジャンルを返す
func RandomGenre() string {
	return GenreOptions[RandomInt(0, len(GenreOptions)-1)]
}

// RandomPurpose はランダムな目的を返す
func RandomPurpose() string {
	return PurposeOptions[RandomInt(0, len(PurposeOptions)-1)]
}

// RandomAuthor はランダムな著者を返す
func RandomAuthor() string {
	return SampleAuthors[RandomInt(0, len(SampleAuthors)-1)]
}

// RandomTitle はランダムなタイトルを返す
func RandomTitle() string {
	return SampleTitles[RandomInt(0, len(SampleTitles)-1)]
}

// ========== Configuration Helpers ==========

// TestConfig はテスト設定を定義
type TestConfig struct {
	EnableIntegrationTests bool
	TestDBPath             string
	LogLevel               string
	Verbose                bool
}

// GetTestConfig はテスト設定を取得
func GetTestConfig() TestConfig {
	return TestConfig{
		EnableIntegrationTests: getEnvBool("RUN_INTEGRATION_TESTS", false),
		TestDBPath:             getEnvString("TEST_DB_PATH", ":memory:"),
		LogLevel:               getEnvString("TEST_LOG_LEVEL", "silent"),
		Verbose:                getEnvBool("TEST_VERBOSE", false),
	}
}

// SkipIntegration は統合テストをスキップするかチェック
func SkipIntegration(t *testing.T) {
	if !GetTestConfig().EnableIntegrationTests {
		t.Skip("Integration tests skipped. Set RUN_INTEGRATION_TESTS=1 to run.")
	}
}

// ========== Test Data Sets ==========

// CreateSampleDataSet はサンプルデータセットを作成
func CreateSampleDataSet() []models.Book {
	return []models.Book{
		{Title: "吾輩は猫である", Author: "夏目漱石", Genre: "Fiction", Purpose: "Entertainment", Description: "猫の視点から描かれた小説"},
		{Title: "Clean Code", Author: "Robert C. Martin", Genre: "Technology", Purpose: "Learning", Description: "Clean code principles"},
		{Title: "The Lean Startup", Author: "Eric Ries", Genre: "Business", Purpose: "Learning", Description: "Startup methodology"},
		{Title: "1984", Author: "George Orwell", Genre: "Fiction", Purpose: "Entertainment", Description: "Dystopian novel"},
		{Title: "Sapiens", Author: "Yuval Noah Harari", Genre: "History", Purpose: "Learning", Description: "A brief history of humankind"},
	}
}

// CreateGenreSpecificDataSet は特定ジャンルのデータセットを作成
func CreateGenreSpecificDataSet(genre string, count int) []models.Book {
	factory := NewBookFactory()
	books := make([]models.Book, count)
	
	for i := 0; i < count; i++ {
		book := factory.CreateBook(WithGenre(genre))
		books[i] = *book
	}
	
	return books
}

// CreatePurposeSpecificDataSet は特定目的のデータセットを作成
func CreatePurposeSpecificDataSet(purpose string, count int) []models.Book {
	factory := NewBookFactory()
	books := make([]models.Book, count)
	
	for i := 0; i < count; i++ {
		book := factory.CreateBook(WithPurpose(purpose))
		books[i] = *book
	}
	
	return books
}

// ========== Helper Functions ==========

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "1" || value == "true" || value == "TRUE"
}

func getEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ========== Pointer Helpers ==========

// StringPtr は文字列のポインタを作成
func StringPtr(s string) *string {
	return &s
}

// UintPtr はuintのポインタを作成
func UintPtr(u uint) *uint {
	return &u
}

// IntPtr はintのポインタを作成
func IntPtr(i int) *int {
	return &i
} 