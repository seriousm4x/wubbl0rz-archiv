package cronjobs

import (
	"github.com/AgileProggers/archiv-backend-go/pkg/logger"
	"github.com/AgileProggers/archiv-backend-go/pkg/router"
	"github.com/robfig/cron"
)

var twitchDownloadsRunning bool

func Init() error {
	c := cron.New()

	logger.Debug.Println("[cronjob] registering job: emote update")
	if err := c.AddFunc("@every 1h", UpdateEmotes); err != nil {
		return err
	}

	logger.Debug.Println("[cronjob] registering job: stream status")
	if err := c.AddFunc("@every 1m", SetStreamStatus); err != nil {
		return err
	}

	logger.Debug.Println("[cronjob] registering job: twitch downloads")
	if err := c.AddFunc("@every 1h", RunTwitchDownloads); err != nil {
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

	if err := DownloadVods(); err != nil {
		logger.Error.Println(err)
	}
	if err := DownloadClips(); err != nil {
		logger.Error.Println(err)
	}
	if err := DownloadGames(); err != nil {
		logger.Error.Println(err)
	}

	twitchDownloadsRunning = false

	// delete cached routes
	if err := router.MemoryStore.Cache.Purge(); err != nil {
		logger.Error.Println(err)
	}
}
