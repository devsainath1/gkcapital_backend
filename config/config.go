package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
	DBHost        string `mapstructure:"db_host"`
	DBPort        string `mapstructure:"db_port"`
	DBUser        string `mapstructure:"db_user"`
	DBPassword    string `mapstructure:"db_password"`
	DBName        string `mapstructure:"db_name"`
	JWTSecret     string `mapstructure:"jwt_secret"`
	ServerPort    string `mapstructure:"server_port"`
	AppEnv        string `mapstructure:"app_env"`
	UploadPath    string `mapstructure:"upload_path"`
	MaxUploadSize string `mapstructure:"max_upload_size"`
}

func LoadConfig() *Config {
	v := viper.New()

	// ── 1. Built-in defaults ─────────────────────────────────────────────────
	v.SetDefault("app_env", "local")
	v.SetDefault("server_port", "8080")
	v.SetDefault("db_host", "localhost")
	v.SetDefault("db_port", "3306")
	v.SetDefault("db_user", "root")
	v.SetDefault("db_password", "")
	v.SetDefault("db_name", "gk_capital")
	v.SetDefault("jwt_secret", "gk-capital-secret-key-change-in-production")
	v.SetDefault("upload_path", "./uploads")
	v.SetDefault("max_upload_size", "10")

	// ── 2. Config file ───────────────────────────────────────────────────────
	// Resolve config file path:
	//   Priority: CONFIG_PATH env var → ./config/local/config.yaml (local dev fallback) → ./config.yaml (fallback)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		if _, err := os.Stat("./config/local/config.yaml"); err == nil {
			configPath = "./config/local/config.yaml"
		} else {
			configPath = "./config.yaml"
		}
	}

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Info: No config file found at '%s'. Using defaults + env vars.", configPath)
		} else {
			log.Printf("Warning: Failed to read config file '%s': %v. Using defaults + env vars.", configPath, err)
		}
	} else {
		log.Printf("Successfully loaded config from '%s'", v.ConfigFileUsed())
	}

	// ── 3. Auto-bind environment variables ───────────────────────────────────
	// DB_HOST → db_host, JWT_SECRET → jwt_secret, etc.
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// ── 4. Unmarshal into Config struct ──────────────────────────────────────
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	log.Printf("Config loaded: env=%s port=%s db=%s@%s:%s/%s",
		cfg.AppEnv, cfg.ServerPort,
		cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

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
