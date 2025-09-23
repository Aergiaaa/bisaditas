package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	_ "github.com/Aergiaaa/bisaditas/docs"
	"github.com/Aergiaaa/bisaditas/internal/database"
	"github.com/Aergiaaa/bisaditas/internal/env"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

// @title Gin Event API
// @version 1.0
// @description This is a sample server for managing events.
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description enter your bearer token in the format **Bearer &lt;token&gt;**

type app struct {
	host      string
	port      int
	jwtSecret string
	models    database.Models
}

func runMigrations() {
	log.Println("Running migrations...")

	err := database.Migrate()
	if err != nil {
		log.Printf("error running migrations: %v", err)
	} else {
		log.Println("Migrations completed successfully.")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	port := env.GetEnvInt("PORT", 8080)
	if herokuPort := os.Getenv("PORT"); herokuPort != "" {
		if p, err := strconv.Atoi(herokuPort); err == nil {
			port = p
		}
	}

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	runMigrations()

	models := database.NewModels(db)
	app := &app{
		host:      env.GetEnvString("HOST", "localhost"),
		port:      port,
		jwtSecret: env.GetEnvString("JWT_SECRET", "secret-123456"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatalf("error serving app: %v", err)
	}
}
