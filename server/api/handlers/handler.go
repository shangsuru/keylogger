package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RecordingsHandler struct {
	db *sql.DB
}

type Recording struct {
	Ip         string    `json:"ip"`
	Timestamp  time.Time `json:"timestamp"`
	Keystrokes string    `json:"keystrokes"`
}

func NewRecordingsHandler(db *sql.DB) *RecordingsHandler {
	return &RecordingsHandler{db: db}
}

func (handler *RecordingsHandler) ListPerDay(c *gin.Context) {
	date := c.Param("day")
	rows, err := handler.db.Query("SELECT * FROM recordings WHERE time_stamp::date = $1::date", date)
	if err != nil {
		log.Printf("Couldn't query the database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()
	var recordings []Recording
	for rows.Next() {
		var recording Recording
		err = rows.Scan(&recording.Ip, &recording.Timestamp, &recording.Keystrokes)
		if err != nil {
			log.Printf("Error reading in rows: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		recordings = append(recordings, recording)
	}

	c.JSON(http.StatusOK, recordings)
}
