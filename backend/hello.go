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

// getBooks responds with the list of all albums as JSON.
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

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)

	router.Run("localhost:8080")
}
