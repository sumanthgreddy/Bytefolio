package handlers

import (
    "net/http"
	"fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gocql/gocql"
    "github.com/sirupsen/logrus"
	"PersonalWebsite/backend/database"
    "PersonalWebsite/backend/middleware" // Adjust the path if necessary
)

type Profile struct {
    Profile     string `json:"profile"` // Matches your primary key 'profile'
    ProfileImage []byte `json:"profile_image"` // For storing the blob  
    Code        string `json:"code"`
    MainSummary string `json:"main_summary"`
}

type WorkExperience struct {
    Profile      string    `json:"profile"`
    CompanyName  string    `json:"company_name"`
    JobTitle     string    `json:"job_title"`
    StartDate    time.Time `json:"start_date"` 
    EndDate      time.Time `json:"end_date"`
    WorkDetails  string    `json:"work_details"`
}


func adminProfileHandler(log *logrus.Logger, session *gocql.Session) gin.HandlerFunc {
    return func(c *gin.Context) {
        profiles, err := database.fetchAllProfiles(session)
        if err != nil {
            log.Errorf("Error fetching profiles: %v", err)
            c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Error fetching profiles. Please try again."})
            return // Important - Stop execution if there's an error
        }

        c.HTML(http.StatusOK, "admin.tmpl", gin.H{"profiles": profiles})
    }
}


func adminLoginHandler(log *logrus.Logger, session *gocql.Session) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Render admin login form 
        c.HTML(http.StatusOK, "admin_login.tmpl", nil)

        // 2. Handle form submission (POST request)
        if c.Request.Method == "POST" {
            email := c.PostForm("email")
            password := c.PostForm("password") 

            // Validate email/password 
            if isValidAdmin(log, session, email, password) {  
                c.Set("isAdmin", true)  
                c.Redirect(http.StatusFound, "/ganjimain99") 
            } else {
                // Show error on the admin_login.tmpl
                c.HTML(http.StatusOK, "error.tmpl", gin.H{"error": "Invalid admin credentials"})
            }
        }
    }
}

// Functions to perform CRUD operations (Placeholders)

func updateProfile(session *gocql.Session, profile Profile) error {
    err := session.Query("UPDATE website.profiles SET full_name = ?, email = ?, phone_number = ? WHERE profile = ?", 
        profile.FullName, profile.Email, profile.PhoneNumber, profile.ProfileID).Exec()
    if err != nil {
        return fmt.Errorf("error updating profile: %v", err)
    }
    return nil 
}

func createProfile(session *gocql.Session, profile Profile) error {
    profile.ProfileID = uuid.New().String()
    err := session.Query("INSERT INTO website.profiles (profile, full_name, email, phone_number) VALUES (?, ?, ?, ?)",
        profile.ProfileID, profile.FullName, profile.Email, profile.PhoneNumber).Exec()
    if err != nil {
        return fmt.Errorf("error creating profile: %v", err)
    }
    return nil 
}

func deleteProfile(session *gocql.Session, profileID string) error {
    err := session.Query("DELETE FROM website.profiles WHERE profile = ?", profileID).Exec()
    if err != nil {
        return fmt.Errorf("error deleting profile: %v", err)
    }
    return nil
}

func isValidAdmin(log *logrus.Logger, session *gocql.Session, email, password string) bool {
    // 1. Fetch admin credentials from database 
    var storedEmail, storedPasswordHash string
    err := session.Query("SELECT email, password_hash FROM admins WHERE email = ?", email).Consistency(gocql.One).Scan(&storedEmail, &storedPasswordHash)
    if err != nil {
        log.Errorf("Error fetching admin data: %v", err)
        return false 
    }

    // 2. Basic password check using bcrypt
    err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
    if err != nil { 
        return false
    }

    // 3. Advanced Authentication (Example: Email Verification)
    // a) Generate and store a random verification code
    verificationCode := generateRandomCode() // You'll need a function for this
    err = session.Query("UPDATE admins SET verification_code = ? WHERE email = ?", verificationCode, email).Exec()
    if err != nil {
        log.Errorf("Error storing verification code: %v", err)
        return false
    }

    // b) Send the code to the admin's email (Implementation not shown here)
    err = sendVerificationEmail(email, verificationCode)
    if err != nil {
        log.Errorf("Error sending verification email: %v", err)
        return false
    }

 // 3. Advanced Authentication 
    verificationCode := generateRandomCode() 
    err = session.Query("UPDATE admins SET verification_code = ? WHERE email = ?", verificationCode, email).Exec()
    if err != nil {
        log.Errorf("Error in generating Random code: %v", err)
		return false
    }

    err = sendVerificationEmail(email, verificationCode)
    if err != nil {
        log.Errorf("Error is sending Verfication email: %v", err)
    }

    return true // Temporary (until code verification is implemented) 
}



func generateRandomCode() string {
    rand.Seed(time.Now().UnixNano())
    chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
    length := 8 
    var b strings.Builder
    for i := 0; i < length; i++ {
        b.WriteRune(chars[rand.Intn(len(chars))])
    }
    return b.String() 
}


func sendVerificationEmail(email, code string) error {
    from := config.Email.Address
    password := config.Email.Password
    to := []string{email}
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    message := []byte(fmt.Sprintf("Your verification code is: %s", code))

    auth := smtp.PlainAuth("", from, password, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
    if err != nil {
        return err
    }
    return nil 
}


func adminCodeVerificationHandler(log *logrus.Logger, session *gocql.Session) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == "POST" {
            email := c.PostForm("email")
            code := c.PostForm("verificationCode")

            // 1. Fetch stored code from database
            var storedCode string
             err := session.Query("SELECT verification_code FROM admins WHERE email = ?", email).Consistency(gocql.One).Scan(&storedCode)
            if err == gocql.ErrNotFound {
                c.HTML(http.StatusOK, "error.tmpl", gin.H{"error": "Email not found"})
                return
            } else if err != nil {
                log.Errorf("Error fetching verification code: %v", err)
                c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "An error occurred"})
                return
            }

            // 2. Compare codes
            if code == storedCode {
				// Clear verification code 
                err = session.Query("UPDATE admins SET verification_code = NULL WHERE email = ?", email).Exec()
                if err != nil {
                    log.Errorf("Error clearing verification code: %v", err)
                    // You might handle this error differently 
                }
                c.Set("isAdmin", true)
                c.Redirect(http.StatusFound, "/ganjimain99")

                c.Set("isAdmin", true)
                c.Redirect(http.StatusFound, "/ganjimain99")
            } else {
                c.HTML(http.StatusOK, "error.tmpl", gin.H{"error": "Invalid code"})
            }
        } else {
            // Render a form to enter code (e.g., admin_code_verification.tmpl)
            c.HTML(http.StatusOK, "admin_code_verification.tmpl", nil)
        }
    }
}

// Cassandra Schema: Ensure your admins table has columns for email, password_hash, and verification_code.
