package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DatabaseUrl    string
	AWSBucket      string
	AWSRegion      string
	AWSAccessKey   string
	AWSSecretKey   string
	Port           string
	JWTSecret      string
	AllowedDomains string
	FileSizeLimit  int64
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if os.Getenv("DATABASE_URL") == "" || os.Getenv("AWS_BUCKET") == "" || os.Getenv("AWS_REGION") == "" || os.Getenv("AWS_ACCESS_KEY") == "" || os.Getenv("AWS_SECRET_KEY") == "" || os.Getenv("JWTSecret") == "" || os.Getenv("ALLOWED_DOMAINS") == "" || os.Getenv("FILE_SIZE_LIMIT") == "" {
		panic("check the .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	file_size_limit, _ := strconv.ParseInt(os.Getenv("FILE_SIZE_LIMIT"), 10, 64)

	return &Config{
		DatabaseUrl:    os.Getenv("DATABASE_URL"),
		AWSAccessKey:   os.Getenv("AWS_ACCESS_KEY"),
		AWSSecretKey:   os.Getenv("AWS_SECRET_KEY"),
		AWSBucket:      os.Getenv("AWS_BUCKET"),
		AWSRegion:      os.Getenv("AWS_REGION"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AllowedDomains: os.Getenv("ALLOWED_DOMAINS"),
		FileSizeLimit:  file_size_limit,
		Port:           port,
	}
}
