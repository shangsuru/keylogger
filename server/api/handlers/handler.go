package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecordingsHandler struct {
	db *sql.DB
}

func NewRecordingsHandler(db *sql.DB) *RecordingsHandler {
	return &RecordingsHandler{db: db}
}

func (handler *RecordingsHandler) ListPerDay(c *gin.Context) {
	c.JSON(http.StatusOK, "ListPerDay")
}

func (handler *RecordingsHandler) SearchForText(c *gin.Context) {
	c.JSON(http.StatusOK, "SearchForText")
}
