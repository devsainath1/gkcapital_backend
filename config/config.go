package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
	DBHost        string `yaml:"db_host"`
	DBPort        string `yaml:"db_port"`
	DBUser        string `yaml:"db_user"`
	DBPassword    string `yaml:"db_password"`
	DBName        string `yaml:"db_name"`
	JWTSecret     string `yaml:"jwt_secret"`
	ServerPort    string `yaml:"server_port"`
	AppEnv        string `yaml:"app_env"`
	UploadPath    string `yaml:"upload_path"`
	MaxUploadSize string `yaml:"max_upload_size"`
}

func LoadConfig() *Config {
	// 1. Determine environment (default to "local" if APP_ENV is empty)
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	// 2. Select folder config/prod or config/local
	var configFolder string
	if env == "prod" || env == "production" {
		configFolder = "config/prod"
	} else {
		configFolder = "config/local"
	}

	// Define defaults
	cfg := &Config{
		AppEnv:        env,
		ServerPort:    "8080",
		DBHost:        "localhost",
		DBPort:        "3306",
		DBUser:        "root",
		DBPassword:    "",
		DBName:        "gk_capital",
		JWTSecret:     "gk-capital-secret-key-change-in-production",
		UploadPath:    "./uploads",
		MaxUploadSize: "10",
	}

	// 3. Load YAML file
	configPath := filepath.Join(configFolder, "config.yaml")
	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("Warning: Failed to open config file at %s (%v). Using defaults and env overrides.", configPath, err)
	} else {
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			log.Printf("Warning: Failed to decode config file %s: %v", configPath, err)
		} else {
			log.Printf("Successfully loaded config from %s", configPath)
		}
	}

	// 4. Override with system environment variables if they are set
	overrideWithEnv := func(key string, target *string) {
		if val := os.Getenv(key); val != "" {
			*target = val
		}
	}

	overrideWithEnv("APP_ENV", &cfg.AppEnv)
	overrideWithEnv("SERVER_PORT", &cfg.ServerPort)
	overrideWithEnv("DB_HOST", &cfg.DBHost)
	overrideWithEnv("DB_PORT", &cfg.DBPort)
	overrideWithEnv("DB_USER", &cfg.DBUser)
	overrideWithEnv("DB_PASSWORD", &cfg.DBPassword)
	overrideWithEnv("DB_NAME", &cfg.DBName)
	overrideWithEnv("JWT_SECRET", &cfg.JWTSecret)
	overrideWithEnv("UPLOAD_PATH", &cfg.UploadPath)
	overrideWithEnv("MAX_UPLOAD_SIZE", &cfg.MaxUploadSize)

	return cfg
}

func ConnectDatabase(cfg *Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	logLevel := logger.Info
	if cfg.AppEnv == "production" || cfg.AppEnv == "prod" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("Database connected successfully")
}
