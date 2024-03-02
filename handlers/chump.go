package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gocql/gocql"
    "github.com/sirupsen/logrus"
    "PersonalWebsite/backend/middleware" // Adjust the path if necessary
)

// SAP Profile Handler
func SapProfileHandler(log *logrus.Logger, session *gocql.Session) gin.HandlerFunc { 
    return func(c *gin.Context) {
        // 1. Authentication 
        userID, err := middleware.AuthenticateUser(c) 
        if err != nil {
            log.Errorf("Error authenticating user: %v", err)
            c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"error": "Unauthorized"})
            return
        }

        // 2. Fetch SAP Profile Data 
        var sapExperience, projects string
        err = session.Query("SELECT sap_exp, projects FROM website.sapprofile WHERE profile = ?", userID).Consistency(gocql.One).Scan(&sapExperience, &projects)
        if err == gocql.ErrNotFound {
            c.HTML(http.StatusNotFound, "error.tmpl", gin.H{"error": "SAP profile not found"})
            return
        } else if err != nil {
            log.Errorf("Error fetching SAP profile: %v", err)
            c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Oops! Something went wrong."})
            return
        }

        // 3. Render the SAP profile template
        c.HTML(http.StatusOK, "error.tmpl", gin.H{
            "sapExperience": sapExperience,
            "projects":      projects, 
            // ... other SAP profile data ...
        })
    }
}