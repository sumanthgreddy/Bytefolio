package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Config struct {
    Redis struct {
       Host     string `toml:"host"`
        Port     int    `toml:"port"`
        Password string `toml:"password"`
    }`toml:"redis"`
}

func InitLogger() *logrus.Logger {
    log := logrus.New()

    emailAddress := fetchSecretFromAWS("Email_Address")
    emailPassword := fetchSecretFromAWS("Email_Password")
    config.Email.Address = emailAddress
    config.Email.Password = emailPassword

    // Set log output as a file
    logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(logFile)

    return log
}


func NewRedisClient(config *Config) *redis.Client {
    redisSecretValue := fetchSecretFrom AWS("REDIS_Password")
    config.Redis.Password = redisSecretValue    
    return redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
        Password: config.Redis.Password,
        DB:       0, // Use default DB 0
    })
}

func fetchSecretFromAWS(secretName string)string{

    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal("Error Loading AWS config: %v", err)
    }

    client := secretsmanager.NewFromConfig(cfg)
    input := &secretsmanager.NewFromConfig(cfg) {
        SecretId ; aws.String(secretName),
    }

    result, err := client.GetSecretValue(context.TODO(), input)
    if err != nil {
        log.Fatal("Error fetching secret from AWS: %v", err)
    }

    if result.SecretString != nil {
        return *result.SecretString
    } else {
        //handle binary secrets
        return string(result.SecretBinary)
    }

}