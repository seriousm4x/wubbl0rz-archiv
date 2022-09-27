package cronjobs

import (
	"github.com/AgileProggers/archiv-backend-go/pkg/logger"
	"github.com/robfig/cron"
)

func Init() error {
	c := cron.New()

	logger.Debug.Println("[cron] registering job: emote update")
	if err := c.AddFunc("@every 1h", UpdateEmotes); err != nil {
		return err
	}

	logger.Debug.Println("[cron] registering job: stream status")
	if err := c.AddFunc("@every 1m", SetStreamStatus); err != nil {
		return err
	}

	logger.Debug.Println("[cron] registering job: twitch downloads")
	if err := c.AddFunc("@every 1h", RunTwitchDownloads); err != nil {
		return err
	}

	c.Start()

	logger.Debug.Printf("[cron] registered %d jobs", len(c.Entries()))

	return nil
}

func RunTwitchDownloads() {
	if err := DownloadVods(); err != nil {
		logger.Error.Println(err)
	}
	if err := DownloadClips(); err != nil {
		logger.Error.Println(err)
	}
	if err := DownloadGames(); err != nil {
		logger.Error.Println(err)
	}
}
