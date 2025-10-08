package database

import (
    "database/sql"
    "log"
    "os"
    _ "github.com/lib/pq" 
)

var DB *sql.DB

func ConnectDB() {
    var err error
    dsn := os.Getenv("DB_DSN")
    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal("Failed to open DB connection:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Failed to ping DB:", err)
    }

    log.Println("Database connection successful")
}