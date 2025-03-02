package main

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var AuthError = errors.New("Authentication failed")

func Authorize(c *gin.Context) error {
	username := c.PostForm("username")
	user, ok := users[username]
	if !ok {
		return AuthError
	}

	// Get session token from cookie
	sessionToken, err := c.Cookie("session_token")
	// print token
	fmt.Println(sessionToken)
	if err != nil || sessionToken == "" || sessionToken != user.SessionToken {
		return AuthError
	}

	// Get CSRF token from header
	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken == "" || csrfToken != user.CSRFToken {
		return AuthError
	}

	return nil
}
