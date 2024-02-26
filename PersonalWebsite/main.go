package main

import (
    "fmt"
    "log"
    "os"

    "github.com/BurntSushi/toml" // Example TOML library
    "github.com/redis/go-redis/v8" // Example Redis library 
    "github.com/gin-gonic/gin"      // Example web framework
	"github.com/sirupsen/logrus" 
	"PersonalWebsite/backend/handlers"
	"github.com/gocql/gocql"
	"github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/keyspaces"
)

type Config struct {
    Redis struct {
        Host     string `toml:"host"`
        Port     int    `toml:"port"`
        Password string `toml:"password"`
    } `toml:"redis"`
    Passwords map[string]string `toml:"passwords"`

	Email struct {
        Address  string `toml:"address"`
        Password string `toml:"password"`
    } `toml:"email"`

}


func main() {
    // Load configuration
    var config Config
    if _, err := toml.DecodeFile("config.toml", &config); err != nil {
        log.Fatal("Error loading config:", err)
    }

	// Initialize logrus
    log := logrus.New()

    // Optionally set log output as a file
    logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer logFile.Close()
    log.SetOutput(logFile)

    // Connect to Redis
    rdb := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
        Password: config.Redis.Password, 
        DB:       0, // Use default DB 0
    })

    // Initialize Gin router (or another web framework of your choice)
    router := gin.Default()

	 setupRoutes(router, log, rdb)

	log.Fatal(router.Run(":8080")) // Assuming you want to run on port 8080

	// Configure AWS client
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatalf("error loading AWS config: %v", err)
    }

    // Connect to Cassandra
    // Connect to Cassandra
    cluster := gocql.NewCluster("cassandra.us-west-2.amazonaws.com:9142") // Use your verified endpoint
    cluster.Keyspace = "website"             
    cluster.Authenticator = gocql.PasswordAuthenticator{
        Username: "sumanth991995@gmail.com",
        Password: "Ganji@99",
    }
    cluster.ProtoVersion = 4 // Ensure this matches your Cassandra version
    session, err := cluster.CreateSession()
    if err != nil {
        log.Fatal("Error connecting to Cassandra:", err)
    }
    defer session.Close() 
}

func setupRoutes(router *gin.Engine, log *logrus.Logger, rdb *redis.Client) {
    router.POST("/login", loginHandler(log, rdb)) 
    router.GET("/personalinfo", authMiddleware(log, rdb), personalInfoHandler(log))
	router.GET("/experience/golang", authMiddleware(log, rdb), golangProfileHandler(log))
	router.GET("/experience/sap", authMiddleware(log, rdb), sapProfileHandler(log))
	router.GET("/experience/ba", authMiddleware(log, rdb), baProfileHandler(log))
	router.GET("/ganjimain99", authMiddleware(log, rdb), baProfileHandler(log))
	router.POST("/ganjimain99", authMiddleware(log, rdb), baProfileHandler(log))
	router.DELETE("/ganjimain99", authMiddleware(log, rdb), baProfileHandler(log))
	router.PUT("/ganjimain99", authMiddleware(log, rdb), baProfileHandler(log))
 
}