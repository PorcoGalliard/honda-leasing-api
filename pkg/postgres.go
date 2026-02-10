package pkg

import (
	"fmt"
	"log"
	"time"

	configs "github.com/nanasuryana335/honda-leasing-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	DB *gorm.DB
}

func InitDB(cfg *configs.Config) (*Postgres, error) {
	dsn := generateDSN(cfg.Database)
	log.Printf("Connecting to database: %s@%s:%s/%s",
		cfg.Database.User,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)

	gormConfig := &gorm.Config{}
	if cfg.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if err := db.Exec("SET search_path TO Honda").Error; err != nil {
		return nil, fmt.Errorf("Failed to set search path: %v ", err)
	}

	log.Printf("✅ Database connected successfully!")
	return &Postgres{DB: db}, nil
}

func generateDSN(dbConfig configs.DatabaseConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
		dbConfig.SSLMode,
		dbConfig.TimeZone)
}

func GetDB(db *Postgres) *gorm.DB {
	if db.DB == nil {
		log.Fatal("Database not initialized. Call InitDB first.")
	}
	return db.DB
}

func CloseDB(db *Postgres) error {
	if db != nil {
		sqlDB, err := db.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
