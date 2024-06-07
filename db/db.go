package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var DB *pgx.Conn

func Init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file: ", err)
	}

	dbUser := os.Getenv("PG_USER")
	dbPassword := os.Getenv("PG_PASSWORD")
	dbName := os.Getenv("PG_DATABASE_NAME")
	dbHost := "localhost" // Имя сервиса базы данных в Docker Compose
	dbPort := os.Getenv("PG_PORT")

	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		logrus.Fatal("Unable to connect to database: ", err)
	}

	DB = conn

	logrus.Info("Database connection established")
}

func GetDB() *pgx.Conn {
	return DB
}
