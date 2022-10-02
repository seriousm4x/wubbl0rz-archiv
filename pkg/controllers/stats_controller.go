package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/database"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"
	"github.com/gin-gonic/gin"
)

type vodPerMonth struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type vodPerWeekday struct {
	Weekday string `json:"weekday"`
	Count   int64  `json:"count"`
}

type startByTime struct {
	Hour  uint  `json:"hour"`
	Count int64 `json:"count"`
}

type longStats struct {
	CountVodsTotal     int64           `json:"count_vods_total"`
	CountClipsTotal    int64           `json:"count_clips_total"`
	CountHoursStreamed float64         `json:"count_h_streamed"`
	CountSizeBytes     uint            `json:"count_size_bytes"`
	TrendVods          int64           `json:"trend_vods"`
	TrendClips         int64           `json:"trend_clips"`
	TrendHoursStreamed float64         `json:"trend_h_streamed"`
	VodsPerMonth       []vodPerMonth   `json:"vods_per_month"`
	VodsPerWeekday     []vodPerWeekday `json:"vods_per_weekday"`
	StartByTime        []startByTime   `json:"start_by_time"`
}

// Get Statistics godoc
// @Summary Get short statistics
// @Tags    Statistics
// @Accept  json
// @Produce  json
// @Success 200 {array} map[string]interface{}
// @Failure  500 {string} string
// @Router  /stats/short [get]
func GetShortStats(c *gin.Context) {
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msg":   "Ok",
		"result": map[string]interface{}{
			"last_vod_sync":   settings.DateVodsUpdate.Format(time.RFC3339),
			"last_emote_sync": settings.DateEmotesUpdate.Format(time.RFC3339),
			"is_live":         settings.IsLive,
		},
	})
}

// Get Statistics godoc
// @Summary Get long statistics
// @Tags    Statistics
// @Accept  json
// @Produce  json
// @Success 200 {array} map[string]interface{}
// @Failure  500 {string} string
// @Router  /stats/long [get]
func GetLongStats(c *gin.Context) {
	var stats longStats
	var vod models.Vod
	var clip models.Clip

	// get CountVodsTotal
	if result := database.DB.Find(&vod).Count(&stats.CountVodsTotal); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get CountClipsTotal
	if result := database.DB.Find(&clip).Count(&stats.CountClipsTotal); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get CountHoursStreamed
	if result := database.DB.Table("vods").Select("sum(duration)/3600").Scan(&stats.CountHoursStreamed); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get CountSizeBytes
	if result := database.DB.Table("vods").Select("sum(size)").Scan(&stats.CountSizeBytes); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get TrendVods
	now := time.Now()
	one_month_ago := now.AddDate(0, -1, 0)
	two_months_ago := now.AddDate(0, -2, 0)
	var count_one_month int64
	var count_two_months int64

	if result := database.DB.Table("vods").Where("date BETWEEN ? AND ?", one_month_ago, now).Count(&count_one_month); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	if result := database.DB.Table("vods").Where("date BETWEEN ? AND ?", two_months_ago, one_month_ago).Count(&count_two_months); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}
	stats.TrendVods = count_one_month - count_two_months

	// get TrendClips
	if result := database.DB.Table("clips").Where("date BETWEEN ? AND ?", one_month_ago, now).Count(&count_one_month); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	if result := database.DB.Table("clips").Where("date BETWEEN ? AND ?", two_months_ago, one_month_ago).Count(&count_two_months); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}
	stats.TrendClips = count_one_month - count_two_months

	// get TrendHoursStreamed
	var count_h_streamed_month1 float64
	var count_h_streamed_month2 float64
	if result := database.DB.Table("vods").Where("date BETWEEN ? AND ?", one_month_ago, now).Select("sum(duration)/3600").Scan(&count_h_streamed_month1); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	if result := database.DB.Table("vods").Where("date BETWEEN ? AND ?", two_months_ago, one_month_ago).Select("sum(duration)/3600").Scan(&count_h_streamed_month2); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}
	stats.TrendHoursStreamed = count_h_streamed_month1 - count_h_streamed_month2

	// get VodsPerMonth
	for i := 11; i >= 0; i-- {
		monthsBack := now.AddDate(0, -i, 0)
		range_start := time.Date(monthsBack.Year(), monthsBack.Month(), 1, 0, 0, 0, 0, monthsBack.Local().Location())
		range_end := range_start.AddDate(0, 1, -1)
		monthStr := fmt.Sprintf("%s %d", range_start.Month().String()[:3], range_start.Year()%100)
		month := vodPerMonth{
			Month: monthStr,
		}

		if result := database.DB.Table("vods").Where("date BETWEEN ? AND ?", range_start, range_end).Count(&month.Count); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msg":   "Failed to get stats",
			})
			return
		}
		stats.VodsPerMonth = append(stats.VodsPerMonth, month)
	}

	// get VodsPerWeekday
	weekdays := []string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}
	var weekday vodPerWeekday
	for i, day := range weekdays {
		if result := database.DB.Table("vods").Where("(extract(dow from date) = ?)", i).Count(&weekday.Count); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msg":   "Failed to get stats",
			})
			return
		}
		weekday.Weekday = day
		stats.VodsPerWeekday = append(stats.VodsPerWeekday, weekday)
	}

	// get StartByTime
	if result := database.DB.Table("vods").Select("extract(hour from date) as hour, count(extract(hour from date)) as count").Group("hour").Order("hour asc").Scan(&stats.StartByTime); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":  false,
		"msg":    "Ok",
		"result": stats,
	})
}
