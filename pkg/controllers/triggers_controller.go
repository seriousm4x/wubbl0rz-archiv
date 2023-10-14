package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/cronjobs"
	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/logger"
)

// Triggers godoc
// @Summary Trigger downloads
// @Tags    Triggers
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Router  /trigger/download [get]
func TriggerDownloads(c *gin.Context) {
	logger.Debug.Println("[trigger] manual download")
	go cronjobs.RunTwitchDownloads()
	c.JSON(http.StatusOK, nil)
}
