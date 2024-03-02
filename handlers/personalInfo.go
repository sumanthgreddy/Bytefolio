package handlers

import (
	"net/http"

	"PersonalWebsite/backend/middleware"
	

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

// personalInfoHandler displays a user's personal information
func PersonalInfoHandler(log *logrus.Logger, session *gocql.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Authentication (replace with your method)

		userID, err := middleware.AuthenticateUser(c) // Replace with your authentication logic
		if err != nil {
			log.Errorf("Error authenticating user: %v", err)
			c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"error": "Unauthorized"})
			return
		}

		// 2. Fetch personal information from Cassandra
		var profile, mainSummary string
		var profileImage []byte
		err = session.Query("SELECT profile, main_summary, profile_image FROM website.main WHERE profile = ?", userID).Consistency(gocql.One).Scan(&profile, &mainSummary, &profileImage)
		if err == gocql.ErrNotFound {
			c.HTML(http.StatusNotFound, "error.tmpl", gin.H{"error": "Profile not found"})
			return
		} else if err != nil {
			log.Errorf("Error fetching profile from Cassandra: %v", err)
			c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Oops! Something went wrong. Please try again later."})
			return
		}
		var icons []byte
		var workExp string
		err = session.Query("SELECT icons, work_exp FROM website.workdetails WHERE profile = ?", userID).Consistency(gocql.One).Scan(&icons, &workExp)
		if err == gocql.ErrNotFound {
			c.HTML(http.StatusNotFound, "error.tmpl", gin.H{"error": "WorkProfile not found"})
			return
		} else if err != nil {
			log.Errorf("Error fetching WorkProfile from Cassandra: %v", err)
			c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Oops! Something went wrong. Please try again later."})
			return
		}

		// 3. Render a template
		c.HTML(http.StatusOK, "personal_info.tmpl", gin.H{
			"name":         name,
			"description":  description,
			"profilePhoto": profilePhoto, // Assuming you handle image rendering in your template
			"icons":        icons,        // Assuming you handle icon rendering in your tempate
		})
	}
}

// authenticateUser moved to middleware> auth.go
