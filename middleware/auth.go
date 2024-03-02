package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func AuthMiddleware(log *logrus.Logger, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.Request.Header.Get("Password") // Adjust how you extract the password if needed

		// Check if password exists as a key in Redis
		val, err := rdb.Get(c.Request.Context(), password).Result()
		if err == redis.Nil {
			log.Warn("Invalid password attempt")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Error("Redis error: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Success! Store the route the user should be sent to
		c.Set("targetRoute", val)
		c.Next()
	}
}

func AuthenticateUser(c *gin.Context, rdb *redis.Client) (string, error) {
	password := c.Request.FormValue("password")

	// Retrieve stored hash from Redis
	storedHash, err := rdb.Get(c.Request.Context(), password).Result()
	if err == redis.Nil {
		// Likely incorrect password
		return "", fmt.Errorf("invalid username or password")
	} else if err != nil {
		// Redis or other error
		return "", fmt.Errorf("error authenticating, please try again")
	}

	// Compare provided password against stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil { // Means password mismatch
		return "", fmt.Errorf("invalid username or password")
	}

	// Success! Retrieve corresponding userID
	userID, err := rdb.Get(c.Request.Context(), "user_id:"+password).Result()
	if err != nil {
		return "", fmt.Errorf("error retrieving user data")
	}

	return userID, nil
}

func AdminAuthMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			log.Warn("Unauthorized access to admin area")
			c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
	}
}
