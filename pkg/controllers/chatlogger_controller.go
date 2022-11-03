package controllers

import (
	"net/http"

	"github.com/AgileProggers/archiv-backend-go/pkg/database"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/gin-gonic/gin"
)

// GetAllChatMessages godoc
// @Summary Get chat messages
// @Tags    ChatMessage
// @Produce  json
// @Success 200 {object} models.ChatMessage
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router  /chat/ [get]
// @Param   from query string true "messages after unix timestamp"
// @Param   to query string false "messages before unix timestamp"
func GetAllChatMessages(c *gin.Context) {
	var messages []models.ChatMessage
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "'from' query param missing",
		})
		return
	}

	query := database.DB.Model(&models.ChatMessage{})

	if toStr != "" {
		query = query.Where("created_at > to_timestamp(?) and created_at < to_timestamp(?)", fromStr, toStr)
	} else {
		query = query.Where("created_at > to_timestamp(?)", fromStr)
	}

	if result := query.Find(&messages); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":  false,
		"msg":    "Ok",
		"result": messages,
	})
}
