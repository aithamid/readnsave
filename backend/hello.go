package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Book represents data about a book.
type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	CurrentPage int    `json:"currentPage"`
	TotalPages  int    `json:"totalPages"`
}

// Books slice to seed record book data.
var books = []Book{
	{ID: "1", Title: "Blue Train", Author: "John Coltrane", CurrentPage: 56, TotalPages: 100},
	{ID: "2", Title: "Jeru", Author: "Gerry Mulligan", CurrentPage: 56, TotalPages: 100},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", CurrentPage: 56, TotalPages: 100},
}

// InitializeDatabase initializes the database with the required tables.
func InitializeDatabase(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// SQL to create tables
		sql := `
		CREATE TABLE IF NOT EXISTS Users (
			userId SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			valid BOOLEAN DEFAULT FALSE,
			public BOOLEAN DEFAULT TRUE
		);

		CREATE TABLE IF NOT EXISTS Books (
			bookId SERIAL PRIMARY KEY,
			bookname VARCHAR(100) NOT NULL,
			pages INT DEFAULT 0,
			totalpages INT NOT NULL,
			userId INT,
			FOREIGN KEY (userId) REFERENCES Users(userId)
		);

		CREATE TABLE IF NOT EXISTS Followers (
			id SERIAL PRIMARY KEY,
			userId1 INT,
			userId2 INT,
			approve BOOLEAN DEFAULT FALSE,
			FOREIGN KEY (userId1) REFERENCES Users(userId),
			FOREIGN KEY (userId2) REFERENCES Users(userId)
		);

		CREATE TABLE IF NOT EXISTS History (
			id SERIAL PRIMARY KEY,
			userId INT,
			bookId INT,
			pagesAdded INT NOT NULL,
			datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (userId) REFERENCES Users(userId),
			FOREIGN KEY (bookId) REFERENCES Books(bookId)
		);`

		// Execute the SQL
		if err := db.Exec(sql).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Database initialized successfully"})
	}
}

func ResetDatabase(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// SQL to drop tables
		sql := `
		DROP TABLE IF EXISTS History;
		DROP TABLE IF EXISTS Followers;
		DROP TABLE IF EXISTS Books;
		DROP TABLE IF EXISTS Users;`

		// Execute the SQL
		if err := db.Exec(sql).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Database reset successfully"})
	}
}

// getBooks responds with the list of all books as JSON.
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id := c.Param("id")
	for _, a := range books {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func postBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func putBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook Book

	if err := c.BindJSON(&updatedBook); err != nil {
		return
	}

	for i, a := range books {
		if a.ID == id {
			books[i] = updatedBook
			c.IndentedJSON(http.StatusOK, updatedBook)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")
	for i, a := range books {
		if a.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func connectDB() *gorm.DB {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Initialize the database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection established")
	return db
}

func main() {
	// Connect to the database
	db := connectDB()

	// Initialize Gin router
	router := gin.Default()

	// Routes
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)
	router.POST("/books", postBook)
	router.PUT("/books/:id", putBook)
	router.DELETE("/books/:id", deleteBook)
	router.POST("/init-db", InitializeDatabase(db)) // Pass db to InitializeDatabase
	router.POST("/reset-db", ResetDatabase(db))     // Pass db to ResetDatabase

	// Start the server
	serverPort := os.Getenv("SERVER_PORT")
	router.Run(":" + serverPort)
}
