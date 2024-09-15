package database

import (
    "log"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "backend/models" // Ensure this import path is correct
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
    // Read database connection details from environment variables
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")
    dbSSLMode := os.Getenv("DB_SSLMODE")
    dbTimeZone := os.Getenv("DB_TIMEZONE")

    if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
        log.Fatal("Database environment variables are not set")
    }

    dsn := "host=" + dbHost +
        " user=" + dbUser +
        " password=" + dbPassword +
        " dbname=" + dbName +
        " port=" + dbPort +
        " sslmode=" + dbSSLMode +
        " TimeZone=" + dbTimeZone

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    log.Println("Database connection established")

    // Migrate the schema
    err = DB.AutoMigrate(&models.User{}, &models.File{}, &models.SharedFile{})
    if err != nil {
        log.Fatalf("Error migrating database schema: %v", err)
    }
    log.Println("Database schema migrated")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
    return DB
}
