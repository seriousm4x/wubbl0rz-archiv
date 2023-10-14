package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/database"
	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/models"
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
	var settings models.Settings
	fromStr := c.Query("from")
	toStr := c.Query("to")

	// check id from and to are empty
	if fromStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "'from' query param missing",
		})
		return
	}
	if toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "'to' query param missing",
		})
		return
	}

	// check id from and to are ints
	fromInt, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}
	toInt, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	// get chat messages
	query := database.DB.Model(&models.ChatMessage{}).Where("created_at > to_timestamp(?) and created_at < to_timestamp(?)", fromInt, toInt)
	if result := query.Find(&messages); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Error while getting chat messages",
		})
		return
	}

	// get broadcaster id
	if result := database.DB.Model(&models.Settings{}).First(&settings); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Error while getting broadcaster id",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":          false,
		"msg":            "Ok",
		"broadcaster_id": settings.BroadcasterId,
		"result":         messages,
	})
}
