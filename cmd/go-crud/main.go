package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joho/godotenv"

	"github.com/cxcnxl/go-crud/internal/routes"
)

func main() {
    parseDotEnv();
    db := connectToDb();
    startServer(db);
}

// prepares .env file so it can be read via os.Getenv or panics
func parseDotEnv() {
    err := godotenv.Load();
    if err != nil {
        slog.Error("Error parsing .env " + err.Error());
        panic(err);
    }
}

// opens connection to the database or panics
func connectToDb() *gorm.DB {
    mysqlConnString := os.Getenv("DB_CONNECTION_STRING");
    if mysqlConnString == "" {
        slog.Error("DB Connection String variable not found");
        panic("DB Connection String variable not found");
    }

    db, err := gorm.Open(
        mysql.Open(mysqlConnString),
        &gorm.Config{},
    );
    if err != nil {
        slog.Error("Error opening connection to the database: " + err.Error());
        panic(err);
    }

    return db;
}

// starts http server or panics
func startServer(db *gorm.DB) {
    router := routes.NewRouter(db);

    const port int = 8080;
    addr := fmt.Sprintf(":%d", port);

    slog.Info(fmt.Sprintf("server is running on localhost:%d", port));

    err := http.ListenAndServe(addr, router.Mux);
    if err != nil {
        slog.Error("Error starting server: " + err.Error());
        panic(err);
    }
}
