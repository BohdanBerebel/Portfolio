package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BohdanBerebel/Portfolio/go_app/internal/auth"
)

func getUserID(c *gin.Context) (int64, bool) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is not authenticated",
		})
		return 0, false
	}

	return userID, true
}
