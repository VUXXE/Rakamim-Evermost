package mysql

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"evermos-api/internal/pkg/entity"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func LoadConfigFromEnv() Config {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "evermos")

	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))
	if maxOpenConns <= 0 {
		maxOpenConns = 100
	}

	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	if maxIdleConns <= 0 {
		maxIdleConns = 10
	}

	connMaxLifetimeStr := getEnv("DB_CONN_MAX_LIFETIME", "30m")
	connMaxLifetime, err := time.ParseDuration(connMaxLifetimeStr)
	if err != nil {
		connMaxLifetime = 30 * time.Minute
	}

	return Config{
		Host:            host,
		Port:            port,
		User:            user,
		Password:        password,
		DBName:          dbName,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}
}

func Connect(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running GORM AutoMigrate...")
	return db.AutoMigrate(
		&entity.User{},
		&entity.Store{},
		&entity.Address{},
		&entity.Category{},
		&entity.Product{},
		&entity.Transaction{},
		&entity.ProductLog{},
	)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
