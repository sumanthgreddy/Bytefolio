package utils

import(
	"math/rand"
	"fmt" 
    "net/smtp" 
    "strings" 
    "time"
    "github.com/golang-jwt/jwt/v4"
	
)

//Generate JWT token
func generateJWT(username string) (string, error) {
    // Create a new token
    token := jwt.New(jwt.SigningMethodHS256)

    // Set claims (payload)
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = username
    claims["exp"] = time.Now().Add(time.Minute * 60).Unix() // Token expiration time (60 min)

    // Sign the token with our secret key
    tokenString, err := token.SignedString(sampleSecretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

//store generated token in Redis
func storeTokenInRedis(token string, username string) error {
    // Your Redis logic here (e.g., set the token with the username as the key)
    // Example: redisClient.Set(username, token, time.Hour*24)

    return nil
}

func GenerateRandomCode() string {
    rand.Seed(time.Now().UnixNano())
    chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
    length := 8 
    var b strings.Builder
    for i := 0; i < length; i++ {
        b.WriteRune(chars[rand.Intn(len(chars))])
    }
    return b.String() 
}

func SendVerificationEmail(email, code string, config Config) error {
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

