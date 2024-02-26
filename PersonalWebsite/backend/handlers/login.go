/*
login.go

Receive password: Extract the password from the submitted form data.
Validation (Optional): Perform basic checks (e.g., not empty, maybe some length requirements).
Redis Lookup: Check if the submitted password exists as a key in Redis.
If found: Redirect to the corresponding page (get the target route from your Redis mapping).
Not found: Display an error message.
*/


package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	
)

// loginHandler handles password submission and redirects if valid
func loginHandler(log *logrus.Logger, rdb *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
		password := c.Request.FormValue("password") // Assuming password is in a form

		//  validation 
		if len(password) < 8 { 
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "Password must be at least 8 characters long"})
			return
		}

		storedHash, err := rdb.Get(c.Request.Context(), password).Result()
        if err == redis.Nil {
            // Likely a wrong password
            c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"error": "Invalid username or password"})
            return
        } else if err != nil {
            // Redis or another internal issue
            c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Error, please try again later"})
            return
        }

        err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
        if err != nil { // Means password mismatch
            c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"error": "Invalid username or password"})
            return
        } 

        // SUCCESS: If regular password is correct 
        userID, err := rdb.Get(c.Request.Context(), "user_id:"+password).Result() 
        if err != nil {
            log.Errorf("Error retrieving user data: %v", err)
            c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Error, please try again later"})
            return
        }

        // Check if it's an admin credential
        isAdmin, err := rdb.Get(c.Request.Context(), "admin:"+userID).Result()
        if err == redis.Nil { 
             // Not an admin, proceed with regular user flow
            c.Redirect(http.StatusFound, "/login")
        } else if err != nil {
            log.Errorf("Redis error: %v", err)
            c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Redis Error, please try again later"})
            return
        } else {
            // Redirect to admin login
            c.Redirect(http.StatusFound, "/ganjimain99/login")
        }
    }
}