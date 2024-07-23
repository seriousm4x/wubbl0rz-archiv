package routes

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"golang.org/x/sync/errgroup"
)

// Route to gather archive statistics
func Stats(app *pocketbase.PocketBase, c echo.Context) error {
	type chatter struct {
		Name     string `json:"name"`
		MsgCount int    `json:"msg_count"`
	}

	stats := struct {
		CountVods  int       `json:"count_vods"`
		CountClips int       `json:"count_clips"`
		CountHours int       `json:"count_hours"`
		CountSize  int       `json:"count_size"`
		TrendVods  int       `json:"trend_vods"`
		TrendClips int       `json:"trend_clips"`
		TrendHours int       `json:"trend_hours"`
		Chatters   []chatter `json:"chatters"`
		LastUpdate time.Time `json:"last_update"`
	}{
		CountVods:  0,
		CountClips: 0,
		CountHours: 0,
		CountSize:  0,
		TrendVods:  0,
		TrendClips: 0,
		TrendHours: 0,
		Chatters:   []chatter{},
		LastUpdate: time.Now().UTC(),
	}

	now := time.Now()
	last30 := now.Add(-(24 * 30) * time.Hour)
	last60to30 := last30.Add(-(24 * 30) * time.Hour)

	errs, _ := errgroup.WithContext(context.Background())
	var wg sync.WaitGroup
	wg.Add(3)

	// process vods
	errs.Go(func() error {
		collection, err := app.Dao().FindCollectionByNameOrId("vod")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "failed to get vods",
			})
		}
		stats.LastUpdate = collection.Updated.Time().UTC()

		vods, err := app.Dao().FindRecordsByExpr("vod", dbx.HashExp{"publish": true})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "failed to get vods",
			})
		}
		stats.CountVods = len(vods)
		countLast30 := 0
		countLast60to30 := 0
		countSecondsLast30 := 0
		countSecondsLast60to30 := 0
		countHoursInSeconds := 0
		for _, vod := range vods {
			duration := vod.GetInt("duration")
			countHoursInSeconds += duration
			stats.CountSize += vod.GetInt("size")
			vodTime := vod.GetDateTime("date").Time()
			if vodTime.After(last30) {
				countLast30++
				countSecondsLast30 += duration
			} else if vodTime.After(last60to30) {
				countLast60to30++
				countSecondsLast60to30 += duration
			}
		}
		stats.CountHours = countHoursInSeconds / 60 / 60
		stats.TrendVods = countLast60to30 - countLast30
		stats.TrendHours = (countSecondsLast60to30 - countSecondsLast30) / 60 / 60
		return nil
	})

	// process clips
	errs.Go(func() error {
		clips, err := app.Dao().FindRecordsByExpr("clip")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "failed to get clips",
			})
		}
		stats.CountClips = len(clips)
		countLast30 := 0
		countLast60to30 := 0
		for _, clip := range clips {
			stats.CountSize += clip.GetInt("size")
			clipTime := clip.GetDateTime("date").Time()
			if clipTime.After(last30) {
				countLast30++
			} else if clipTime.After(last60to30) {
				countLast60to30++
			}
		}
		stats.TrendClips = countLast60to30 - countLast30
		return nil
	})

	// process chatmessages

	errs.Go(func() error {
		err := app.Dao().DB().NewQuery("select chatmessage.user_name as name, count(chatmessage.id) as msg_count from chatmessage where chatmessage.user_name not in ('nightbot', 'moobot', 'streamlabs', 'streamelements', 'wizebot', 'deepbot', 'coebot', 'phantombot', 'stay_hydrated_bot') group by chatmessage.user_name order by msg_count desc limit 8").All(&stats.Chatters)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "failed to get chatmessages",
			})
		}
		return nil
	})

	if err := errs.Wait(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, stats)
}
