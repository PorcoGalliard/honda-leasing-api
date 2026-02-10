package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"ENVIRONMENT"`
	Server      ServerConfig   `mapstructure:"SERVER"`
	Database    DatabaseConfig `mapstructure:"DATABASE"`
	JWT         JWTConfig      `mapstructure:"JWT"`
	Storage     StorageConfig  `mapstructure:"STORAGE"`
	CORS        CORSConfig     `mapstructure:"CORS"`
}
type ServerConfig struct {
	Address        string   `mapstructure:"ADDRESS"`
	TrustedProxies []string `mapstructure:"TRUSTED_PROXIES"`
	ReadTimeout    int      `mapstructure:"READ_TIMEOUT"`
	WriteTimeout   int      `mapstructure:"WRITE_TIMEOUT"`
}
type DatabaseConfig struct {
	Host            string `mapstructure:"HOST"`
	Port            string `mapstructure:"PORT"`
	User            string `mapstructure:"USER"`
	Password        string `mapstructure:"PASSWORD"`
	Name            string `mapstructure:"NAME"`
	SSLMode         string `mapstructure:"SSL_MODE"`
	TimeZone        string `mapstructure:"TIMEZONE"`
	MaxOpenConns    int    `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns    int    `mapstructure:"MAX_IDLE_CONNS"`
	ConnMaxLifetime int    `mapstructure:"CONN_MAX_LIFETIME"`
}
type JWTConfig struct {
	Secret string `mapstructure:"SECRET"`
	Expiry int    `mapstructure:"EXPIRY_HOURS"`
}
type StorageConfig struct {
	BasePath     string              `mapstructure:"BASE_PATH"`
	MaxFileSize  int64               `mapstructure:"MAX_FILE_SIZE"`
	PublicURL    string              `mapstructure:"PUBLIC_URL"`
	AllowedTypes AllowedTypesConfig  `mapstructure:"ALLOWED_TYPES"`
	Employees    ModuleStorageConfig `mapstructure:"EMPLOYEES"`
	Products     ModuleStorageConfig `mapstructure:"PRODUCTS"`
	Documents    ModuleStorageConfig `mapstructure:"DOCUMENTS"`
}
type AllowedTypesConfig struct {
	Image    []string `mapstructure:"IMAGE"`
	Document []string `mapstructure:"DOCUMENT"`
}
type ModuleStorageConfig struct {
	MaxSize      int64    `mapstructure:"MAX_SIZE"`
	AllowedTypes []string `mapstructure:"ALLOWED_TYPES"`
	Subdirectory string   `mapstructure:"SUBDIRECTORY"`
}
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"ALLOWED_ORIGINS"`
	AllowedMethods   []string `mapstructure:"ALLOWED_METHODS"`
	AllowedHeaders   []string `mapstructure:"ALLOWED_HEADERS"`
	AllowCredentials bool     `mapstructure:"ALLOW_CREDENTIALS"`
}

func Load() *Config {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	configDir, err := getConfigDir()
	if err != nil {
		log.Printf("Warning: Directory config not found: %v", err)
		configDir = "."
	}

	configPath := filepath.Join(configDir, "configs."+env+".toml")
	log.Printf("File config founded at: %s", configPath)
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	setDefaults()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file %s: %v, using defaults",
			configPath, err)
	} else {
		log.Printf("Loaded configuration from: %s", viper.ConfigFileUsed())
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Error unmarshaling config: %v", err)
	}
	return &config
}

func getConfigDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("could not get current file path")
	}
	return filepath.Dir(filename), nil
}

func setDefaults() {
	// environment
	viper.SetDefault("ENVIRONMENT", "development")
	// Server defaults
	viper.SetDefault("SERVER.ADDRESS", ":8080")
	// Database defaults
	viper.SetDefault("DATABASE.HOST", "localhost")
	viper.SetDefault("DATABASE.PORT", "5434")
	viper.SetDefault("DATABASE.USER", "suryana")
	viper.SetDefault("DATABASE.PASSWORD", "suryana321")
	viper.SetDefault("DATABASE.NAME", "leasing_db")
	viper.SetDefault("DATABASE.SSL_MODE", "disable")
	viper.SetDefault("DATABASE.TIMEZONE", "Asia/Jakarta")
	viper.SetDefault("DATABASE.MAX_OPEN_CONNS", 25)
	viper.SetDefault("DATABASE.MAX_IDLE_CONNS", 5)
	viper.SetDefault("DATABASE.CONN_MAX_LIFETIME", 300)
	// JWT defaults
	viper.SetDefault("JWT.SECRET", "my-secret")
	viper.SetDefault("JWT.EXPIRY_HOURS", 24)
	// Storage defaults
	viper.SetDefault("STORAGE.BASE_PATH", "./uploads")
	viper.SetDefault("STORAGE.MAX_FILE_SIZE", 10485760) // 10MB
	viper.SetDefault("STORAGE.PUBLIC_URL", "/uploads/")
	// Storage allowed types defaults
	viper.SetDefault("STORAGE.ALLOWED_TYPES.IMAGE", []string{"image/jpeg",
		"image/jpg", "image/png", "image/gif"})
	viper.SetDefault("STORAGE.ALLOWED_TYPES.DOCUMENT", []string{"application/pdf",
		"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"})
	// Employee storage defaults
	viper.SetDefault("STORAGE.EMPLOYEES.MAX_SIZE", 5242880) // 5MB
	viper.SetDefault("STORAGE.EMPLOYEES.ALLOWED_TYPES", []string{"image/jpeg",
		"image/jpg", "image/png"})
	viper.SetDefault("STORAGE.EMPLOYEES.SUBDIRECTORY", "employees")
	// CORS defaults
	viper.SetDefault("CORS.ALLOWED_ORIGINS", []string{"http://localhost:3000",
		"http://127.0.0.1:3000"})
	viper.SetDefault("CORS.ALLOWED_METHODS", []string{"GET", "POST", "PUT",
		"DELETE", "OPTIONS"})
	viper.SetDefault("CORS.ALLOWED_HEADERS", []string{"Content-Type",
		"Authorization"})
	viper.SetDefault("CORS.ALLOW_CREDENTIALS", true)
	viper.SetDefault("SERVER.TRUSTED_PROXIES", []string{"127.0.0.1", "localhost",
		"::1"})
}
