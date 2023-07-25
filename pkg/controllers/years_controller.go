package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/database"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/models"
)

// GetYears godoc
// @Summary Get count of vods per year
// @Tags    Years
// @Produce  json
// @Success 200 {object} []map[string]interface{}
// @Failure 400 {string} string
// @Router  /years/ [get]
func GetYears(c *gin.Context) {
	type Year struct {
		Year  int `json:"year"`
		Count int `json:"count"`
	}
	var years []Year
	var vod models.Vod
	database.DB.Model(&vod).Select("to_char(date, 'yyyy') AS year, COUNT(to_char(date, 'yyyy')) AS count").Group("to_char(date, 'yyyy')").Order("year DESC").Find(&years)
	if len(years) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No years found",
			"result":  years,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":  false,
		"msg":    "Ok",
		"result": years,
	})
}
