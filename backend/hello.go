package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// book represents data about a book.
type book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	CurrentPage int    `json:"currentPage"`
	TotalPages  int    `json:"totalPages"`
}

// books slice to seed record book data.
var books = []book{
	{ID: "1", Title: "Blue Train", Author: "John Coltrane", CurrentPage: 56, TotalPages: 100},
	{ID: "2", Title: "Jeru", Author: "Gerry Mulligan", CurrentPage: 56, TotalPages: 100},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", CurrentPage: 56, TotalPages: 100},
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
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func putBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook book

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

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)
	router.POST("/books", postBook)
	router.PUT("/books/:id", putBook)
	router.DELETE("/books/:id", deleteBook)

	router.Run("localhost:8080")
}
