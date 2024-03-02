package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gocql/gocql"
    "github.com/sirupsen/logrus"
    "PersonalWebsite/backend/middleware" // Adjust the path if necessary
)

// Bach Profile Handler
func BachProfileHandler(log *logrus.Logger, session *gocql.Session) gin.HandlerFunc { 
    return func(c *gin.Context) {
        // 1. Authentication 
        userID, err := middleware.AuthenticateUser(c, rdb) 
        if err != nil {
            log.Errorf("Error authenticating user: %v", err)
            c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"error": "Unauthorized"})
            return
        }

        // 2. Fetch Bach Profile Data 
        var bachInterest, compositions string 
        err = session.Query("SELECT bach_interest, compositions FROM website.bachprofile WHERE profile = ?", userID).Consistency(gocql.One).Scan(&bachInterest, &compositions)
        if err == gocql.ErrNotFound {
            c.HTML(http.StatusNotFound, "error.tmpl", gin.H{"error": "Bach profile not found"})
            return
        } else if err != nil {
            log.Errorf("Error fetching Bach profile: %v", err)
            c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Oops! Something went wrong."})
            return
        }

        // 3. Render the Bach profile template
        c.HTML(http.StatusOK, "error.tmpl", gin.H{
            "bachInterest": bachInterest,
            "compositions": compositions, 
        })
    }
}