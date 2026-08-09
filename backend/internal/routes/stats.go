package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/sync/errgroup"
)

// Route to gather archive statistics
func Stats(app core.App, e *core.RequestEvent) error {
	type chatter struct {
		Name     string `json:"name"`
		MsgCount int    `json:"msg_count"`
	}

	type vodStats struct {
		CountVods          int `db:"count_vods"`
		CountLast30        int `db:"count_last30"`
		CountLast60to30    int `db:"count_last60to30"`
		TotalDuration      int `db:"total_duration"`
		DurationLast30     int `db:"duration_last30"`
		DurationLast60to30 int `db:"duration_last60to30"`
		TotalSize          int `db:"total_size"`
	}

	type clipStats struct {
		CountClips      int `db:"count_clips"`
		CountLast30     int `db:"count_last30"`
		CountLast60to30 int `db:"count_last60to30"`
		TotalSize       int `db:"total_size"`
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
		Chatters:   []chatter{},
		LastUpdate: time.Now().UTC(),
	}

	now := time.Now().UTC()
	last30 := now.Add(-30 * 24 * time.Hour)
	last60to30 := last30.Add(-30 * 24 * time.Hour)

	// PocketBase stores dates as strings in SQLite, so format them
	// in the same format used by PocketBase's datetime fields.
	last30Str := last30.Format("2006-01-02 15:04:05.000Z")
	last60to30Str := last60to30.Format("2006-01-02 15:04:05.000Z")

	errs, _ := errgroup.WithContext(context.Background())

	// Process VOD statistics.
	errs.Go(func() error {
		collection, err := app.FindCollectionByNameOrId("vod")
		if err != nil {
			return e.JSON(http.StatusInternalServerError, map[string]any{
				"message": "failed to get vod collection",
			})
		}

		stats.LastUpdate = collection.Updated.Time().UTC()

		var result vodStats

		err = app.DB().NewQuery(`
			SELECT
				COUNT(*) AS count_vods,

				COALESCE(SUM(duration), 0) AS total_duration,

				COALESCE(SUM(size), 0)
					+ COALESCE(SUM(size_audio), 0) AS total_size,

				COALESCE(SUM(
					CASE
						WHEN date >= {:last30} THEN 1
						ELSE 0
					END
				), 0) AS count_last30,

				COALESCE(SUM(
					CASE
						WHEN date >= {:last60to30}
							AND date < {:last30}
						THEN 1
						ELSE 0
					END
				), 0) AS count_last60to30,

				COALESCE(SUM(
					CASE
						WHEN date >= {:last30} THEN duration
						ELSE 0
					END
				), 0) AS duration_last30,

				COALESCE(SUM(
					CASE
						WHEN date >= {:last60to30}
							AND date < {:last30}
						THEN duration
						ELSE 0
					END
				), 0) AS duration_last60to30

			FROM vod
			WHERE publish = true
		`).Bind(dbx.Params{
			"last30":     last30Str,
			"last60to30": last60to30Str,
		}).One(&result)

		if err != nil {
			return e.JSON(http.StatusInternalServerError, map[string]any{
				"message": "failed to get vod statistics",
			})
		}

		stats.CountVods = result.CountVods
		stats.CountHours = result.TotalDuration / 60 / 60
		stats.CountSize += result.TotalSize

		stats.TrendVods = result.CountLast60to30 - result.CountLast30
		stats.TrendHours =
			(result.DurationLast60to30 - result.DurationLast30) / 60 / 60

		return nil
	})

	// Process clip statistics.
	errs.Go(func() error {
		var result clipStats

		err := app.DB().NewQuery(`
			SELECT
				COUNT(*) AS count_clips,

				COALESCE(SUM(size), 0) AS total_size,

				COALESCE(SUM(
					CASE
						WHEN date >= {:last30} THEN 1
						ELSE 0
					END
				), 0) AS count_last30,

				COALESCE(SUM(
					CASE
						WHEN date >= {:last60to30}
							AND date < {:last30}
						THEN 1
						ELSE 0
					END
				), 0) AS count_last60to30

			FROM clip
		`).Bind(dbx.Params{
			"last30":     last30Str,
			"last60to30": last60to30Str,
		}).One(&result)

		if err != nil {
			return e.JSON(http.StatusInternalServerError, map[string]any{
				"message": "failed to get clip statistics",
			})
		}

		stats.CountClips = result.CountClips
		stats.CountSize += result.TotalSize
		stats.TrendClips = result.CountLast60to30 - result.CountLast30

		return nil
	})

	// Process chat messages.
	errs.Go(func() error {
		err := app.DB().NewQuery(`
			SELECT
				user_name AS name,
				COUNT(*) AS msg_count
			FROM chatmessage
			WHERE user_name NOT IN (
				'nightbot',
				'moobot',
				'streamlabs',
				'streamelements',
				'wizebot',
				'deepbot',
				'coebot',
				'phantombot',
				'stay_hydrated_bot'
			)
			GROUP BY user_name
			ORDER BY msg_count DESC
			LIMIT 8
		`).All(&stats.Chatters)

		if err != nil {
			return e.JSON(http.StatusInternalServerError, map[string]any{
				"message": "failed to get chatmessage statistics",
			})
		}

		return nil
	})

	if err := errs.Wait(); err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, stats)
}
