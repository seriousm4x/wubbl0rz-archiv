package controllers

import (
	"fmt"
	"net/http"
	"os"
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

type clipsPerCreator struct {
	Name      string `json:"name"`
	ClipCount int64  `json:"clip_count"`
	ViewCount int64  `json:"view_count"`
}

type longStats struct {
	CountVodsTotal       int64             `json:"count_vods_total"`
	CountClipsTotal      int64             `json:"count_clips_total"`
	CountHoursStreamed   float64           `json:"count_h_streamed"`
	CountSizeBytes       uint              `json:"count_size_bytes"`
	CountTranscriptWords int64             `json:"count_transcript_words"`
	CountUniqueWords     int64             `json:"count_unique_words"`
	CountAvgWords        float64           `json:"count_avg_words"`
	TrendVods            int64             `json:"trend_vods"`
	TrendClips           int64             `json:"trend_clips"`
	TrendHoursStreamed   float64           `json:"trend_h_streamed"`
	DatabaseSize         int64             `json:"database_size"`
	VodsPerMonth         []vodPerMonth     `json:"vods_per_month"`
	VodsPerWeekday       []vodPerWeekday   `json:"vods_per_weekday"`
	StartByTime          []startByTime     `json:"start_by_time"`
	ClipsByCreator       []clipsPerCreator `json:"clips_per_creator"`
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
	if result := database.DB.Model(&vod).Count(&stats.CountVodsTotal); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get CountClipsTotal
	if result := database.DB.Model(&clip).Count(&stats.CountClipsTotal); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get CountHoursStreamed
	if result := database.DB.Model(&vod).Select("sum(duration)/3600").Scan(&stats.CountHoursStreamed); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get CountSizeBytes
	if result := database.DB.Model(&vod).Select("sum(size)").Scan(&stats.CountSizeBytes); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get CountTranscriptWords, CountUniqueWords and CountAvgWords
	tempDest := struct {
		CountTranscriptWords int64   `json:"count_transcript_words"`
		CountUniqueWords     int64   `json:"count_unique_words"`
		CountAvgWords        float64 `json:"count_avg_words"`
	}{}
	if result := database.DB.Raw("select sum(stats.nentry) as count_transcript_words, count(stats.word) as count_unique_words, sum(stats.nentry)/? as count_avg_words from ts_stat('select vods.transcript_vector from vods where vods.publish = true and vods.transcript is not null') as stats", stats.CountVodsTotal).Scan(&tempDest); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}
	stats.CountTranscriptWords = tempDest.CountTranscriptWords
	stats.CountUniqueWords = tempDest.CountUniqueWords
	stats.CountAvgWords = tempDest.CountAvgWords

	// get TrendVods
	if result := database.DB.Raw("select (select count(vods.uuid) from vods where vods.date between (now() - interval '1 month') and now()) - count(vods.uuid) as vods_trend from vods where vods.date between (now() - interval '2 month') and (now() - interval '1 month')").Scan(&stats.TrendVods); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get TrendClips
	if result := database.DB.Raw("select (select count(clips.uuid) from clips where clips.date between (now() - interval '1 month') and now()) - count(clips.uuid) as clips_trend from clips where clips.date between (now() - interval '2 month') and (now() - interval '1 month')").Scan(&stats.TrendClips); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get TrendHoursStreamed
	if result := database.DB.Raw("select (select sum(vods.duration)/3600 from vods where vods.date between (now() - interval '1 month') and now()) - sum(vods.duration)/3600 as trend_h_streamed from vods where vods.date between (now() - interval '2 month') and (now() - interval '1 month')").Scan(&stats.TrendHoursStreamed); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get DatabaseSize
	if result := database.DB.Raw("select pg_database_size(?) as database_size", os.Getenv("POSTGRES_DB")).Find(&stats.DatabaseSize); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get VodsPerMonth
	now := time.Now()
	for i := 11; i >= 0; i-- {
		monthsBack := now.AddDate(0, -i, 0)
		range_start := time.Date(monthsBack.Year(), monthsBack.Month(), 1, 0, 0, 0, 0, monthsBack.Local().Location())
		range_end := range_start.AddDate(0, 1, -1)
		monthStr := fmt.Sprintf("%s %d", range_start.Month().String()[:3], range_start.Year()%100)
		month := vodPerMonth{
			Month: monthStr,
		}

		if result := database.DB.Model(&vod).Where("date BETWEEN ? AND ?", range_start, range_end).Count(&month.Count); result.Error != nil {
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
		if result := database.DB.Model(&vod).Where("(extract(dow from date) = ?)", i).Count(&weekday.Count); result.Error != nil {
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
	if result := database.DB.Model(&vod).Select("extract(hour from date) as hour, count(extract(hour from date)) as count").Group("hour").Order("hour asc").Scan(&stats.StartByTime); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get top clip creators
	if result := database.DB.Model(&clip).Select("creators.name, count(clips.uuid) as clip_count, sum(clips.viewcount) as view_count").Joins("left join creators on clips.creator_uuid = creators.uuid").Group("clips.creator_uuid, creators.name").Order("view_count desc").Limit(15).Scan(&stats.ClipsByCreator); result.Error != nil {
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
