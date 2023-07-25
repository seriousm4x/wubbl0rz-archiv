package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/database"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/models"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/queries"
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

type messagesPerUser struct {
	Name     string `json:"name"`
	MsgCount int64  `json:"msg_count"`
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
	CountChatMessages    int64             `json:"count_chat_messages"`
	VodsPerMonth         []vodPerMonth     `json:"vods_per_month"`
	VodsPerWeekday       []vodPerWeekday   `json:"vods_per_weekday"`
	StartByTime          []startByTime     `json:"start_by_time"`
	ClipsByCreator       []clipsPerCreator `json:"clips_per_creator"`
	MessagesByUser       []messagesPerUser `json:"messages_per_user"`
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

	// get CountChatMessages
	if result := database.DB.Raw("select count(chat_messages.id) as count_chat_messages from chat_messages").Find(&stats.CountChatMessages); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get VodsPerMonth
	if result := database.DB.Raw("select month || ' ' || year as month, count from (select to_char(vods.date, 'Mon') as month, to_char(vods.date, 'MM') as month_int, to_char(vods.date, 'YY') as year, count(vods.uuid) as count from vods where publish = true group by year, month_int, month order by year desc, month_int desc limit 12) as months order by months.year, months.month_int").Scan(&stats.VodsPerMonth); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
	}

	// get VodsPerWeekday
	if result := database.DB.Raw("select weekday, count(dow) as count from (select to_char(vods.date, 'Day') as weekday, extract(dow from vods.date) as dow_int from vods) as dow group by weekday, dow_int order by dow.dow_int").Scan(&stats.VodsPerWeekday); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to get stats",
		})
		return
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

	// top chatter
	if result := database.DB.Raw("select chat_messages.user_display_name as name, count(chat_messages.id) as msg_count from chat_messages where chat_messages.user_name not in ('nightbot', 'moobot', 'streamlabs', 'streamelements', 'wizebot', 'deepbot', 'coebot', 'phantombot', 'stay_hydrated_bot') group by chat_messages.user_display_name order by msg_count desc limit 15").Scan(&stats.MessagesByUser); result.Error != nil {
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
