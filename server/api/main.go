package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/shangsuru/keylogger/server/api/handlers"
)

var recordingsHandler *handlers.RecordingsHandler

func init() {
	// Connect to Postgres
	db, err := sql.Open("postgres", os.Getenv("PSQL_CONN"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	log.Println("Connected to database")

	// Initialize handler
	recordingsHandler = handlers.NewRecordingsHandler(db)
}

func main() {
	router := gin.Default()
	router.GET("/api/recordings/:year/:month/:day", recordingsHandler.ListPerDay)
	router.GET("/api/recordings/search", recordingsHandler.SearchForText)
	router.Run(":5000")
}
