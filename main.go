package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	router "personaData/api"
	db "personaData/db"
)

func main() {
	logrus.Info("Initializing database connection")
	db.Init()

	logrus.Info("Setting up router")
	router := router.SetupRouter()

	logrus.Info("Starting server on :8080")
	logrus.Fatal(http.ListenAndServe(":8080", router))
}
