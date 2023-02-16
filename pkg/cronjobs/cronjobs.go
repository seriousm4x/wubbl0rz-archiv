package cronjobs

import (
	"github.com/AgileProggers/archiv-backend-go/pkg/logger"
	"github.com/AgileProggers/archiv-backend-go/pkg/router"
	"github.com/robfig/cron/v3"
)

var twitchDownloadsRunning bool

func Init() error {
	c := cron.New()

	logger.Debug.Println("[cronjob] registering job: emote update")
	if _, err := c.AddFunc("@every 1h", UpdateEmotes); err != nil {
		return err
	}

	logger.Debug.Println("[cronjob] registering job: stream status")
	if _, err := c.AddFunc("@every 1m", SetStreamStatus); err != nil {
		return err
	}

	logger.Debug.Println("[cronjob] registering job: twitch downloads")
	if _, err := c.AddFunc("@every 1h", RunTwitchDownloads); err != nil {
		return err
	}

	c.Start()

	logger.Debug.Printf("[cronjob] registered %d jobs", len(c.Entries()))

	return nil
}

func RunTwitchDownloads() {
	if twitchDownloadsRunning {
		return
	}

	// run downloads
	twitchDownloadsRunning = true
	downloaded_items := 0

	count, err := DownloadVods()
	if err != nil {
		logger.Error.Println(err)
	}
	downloaded_items += count

	count, err = DownloadClips()
	if err != nil {
		logger.Error.Println(err)
	}
	downloaded_items += count

	if err := DownloadGames(); err != nil {
		logger.Error.Println(err)
	}

	twitchDownloadsRunning = false

	// delete cached routes
	if downloaded_items > 0 {
		if err := router.MemoryStore.Cache.Purge(); err != nil {
			logger.Error.Println(err)
		}
	}
}
