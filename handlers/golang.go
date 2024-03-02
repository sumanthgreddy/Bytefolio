package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gocql/gocql"
    "github.com/sirupsen/logrus"
    "PersonalWebsite/backend/middleware"
)

// Golang Profile Handler
func GolangProfileHandler(log *logrus.Logger, session *gocql.Session) gin.HandlerFunc { 
    return func(c *gin.Context) {
        // 1. Authentication 
        userID, err := middleware.AuthenticateUser(c) 
        if err != nil {
            log.Errorf("Error authenticating user: %v", err)
            c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"error": "Unauthorized"})
            return
        }

        // 2. Fetch Golang Profile Data 
        var golangExperience, projects string
        err = session.Query("SELECT golang_exp, projects FROM website.golangprofile WHERE profile = ?", userID).Consistency(gocql.One).Scan(&golangExperience, &projects)
        if err == gocql.ErrNotFound {
            c.HTML(http.StatusNotFound, "error.tmpl", gin.H{"error": "Golang profile not found"})
            return
        } else if err != nil {
            log.Errorf("Error fetching Golang profile: %v", err)
            c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Oops! Something went wrong."})
            return
        }

        // 3. Render the Golang profile template
        c.HTML(http.StatusOK, "error.tmpl", gin.H{
            "golangExperience": golangExperience,
            "projects":         projects, 
        })
    }
}
