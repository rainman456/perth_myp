package config

import (
	//"log"
	//"net/url"
	"os"
	"strconv"
	//"time"
)

type Config struct {
	RedisAddr string
	RedisPass string
	RedisDB   int
	// Other fields...
	PaystackSecretKey   string
	PaystackPublicKey   string
	PlatformCommission  float64
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	// Email configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func Load() *Config {
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	// AccessTokenExp, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP_MINUTES")) // e.g., 15
	// RefreshTokenExp, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP_DAYS"))  // e.g., 7
	// AccessTokenExp = time.Duration(AccessTokenExp) * time.Minute
	// RefreshTokenExp = time.Duration(RefreshTokenExp) * 24 * time.Hour
	commission, _ := strconv.ParseFloat(os.Getenv("PLATFORM_COMMISSION"), 64)
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return &Config{
		RedisAddr: os.Getenv("REDIS_ADDR"), // e.g., "localhost:6379"
		RedisPass: os.Getenv("REDIS_PASS"),
		RedisDB:   redisDB, // Default 0
		// ...
		PaystackSecretKey:   os.Getenv("PAYSTACK_SECRET_KEY"),
		PaystackPublicKey:   os.Getenv("PAYSTACK_PUBLIC_KEY"),
		PlatformCommission:  commission,
		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
		// Email configuration
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     smtpPort,
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
	}
}
