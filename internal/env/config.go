package env

import (
	"fmt"
	"time"
)

type Config struct {
	Server   ServerCfg
	Database DatabaseCfg
	Storage  StorageCfg
	Security SecurityCfg
}

type ServerCfg struct {
	Port           string
	AllowedDomains []string
	FileSizeLimit  int64
	RateLimit      RateLimitCfg
}

type RateLimitCfg struct {
	Enabled           bool
	RequestsPerMinute int
	BurstSize         int
}

type DatabaseCfg struct {
	Url         string
	MaxOpenCons int
	MaxIdleCons int
	MaxIdleTime time.Duration
}

type SecurityCfg struct {
	JWTSecret      string
	JWTExpiryHours int
	AUD            string
	BcryptCost     int
	AllowedOrigins []string
}

type StorageCfg struct {
	Provider      string
	S3Region      string
	S3Bucket      string
	S3AccessKey   string
	S3SecretKey   string
	LocalBasePath string
}

func NewConfig() *Config {
	return &Config{
		Server: ServerCfg{
			Port:          GetString("SERVER_PORT", ":8080"),
			FileSizeLimit: int64(GetInt("FILE_SIZE_LIMIT", 10*1024*1024)),
			RateLimit: RateLimitCfg{
				Enabled:           GetBool("ENABLE_RATE_LIMIT", true),
				RequestsPerMinute: GetInt("RATE_LIMIT_RPM", 60),
				BurstSize:         GetInt("RATE_LIMIT_BURST", 10),
			},
		},
		Database: DatabaseCfg{
			Url:         GetString("DB_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
			MaxOpenCons: GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleCons: GetInt("DB_MAX_IDLE_CONNS", 25),
			//MaxIdleTime: GetString("DB_MAX_IDLE_TIME", "30"),
			MaxIdleTime: time.Duration(GetInt("DB_MAX_IDLE_TIME_SECONDS", 30)) * time.Second,
		},
		Storage: StorageCfg{
			Provider:      GetString("STORAGE_PROVIDER", "local"),
			S3Region:      GetString("AWS_REGION", ""),
			S3Bucket:      GetString("AWS_BUCKET", ""),
			S3AccessKey:   GetString("AWS_ACCESS_KEY", ""),
			S3SecretKey:   GetString("AWS_SECRET_KEY", ""),
			LocalBasePath: GetString("LOCAL_STORAGE_PATH", "./storage"),
		},
		Security: SecurityCfg{
			JWTSecret:      GetString("JWT_SECRET", "your-secret-key"),
			AUD:            GetString("AUD", "cloudnest"),
			JWTExpiryHours: GetInt("JWT_EXPIRY_HOURS", 24),
			BcryptCost:     GetInt("BCRYPT_COST", 12),
			AllowedOrigins: GetArray("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		},
	}
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Database.Url == "" {
		return fmt.Errorf("database URL is required")
	}
	if c.Security.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	return nil
}
