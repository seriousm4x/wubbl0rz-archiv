package hooks

import (
	"os"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/external"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/cronjobs"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

func InitBackend(app *pocketbase.PocketBase) error {
	// import env to database
	if err := ImportEnv(app); err != nil {
		return err
	}

	// update bearer token
	if err := external.TwitchUpdateBearer(app); err != nil {
		return err
	}

	// update broadcaster id
	if err := external.TwitchUpdateBroadcasterId(app); err != nil {
		return err
	}

	// update run chatlogger
	cl, err := NewChatlogger(app)
	if err != nil {
		return err
	}
	go cl.Run(os.Getenv("BROADCASTER_NAME"))

	// run discord bot
	go RunDiscordBot(app)

	// run all cronjobs once
	cronjobs.SetStreamStatus(app)
	publicSettings, err := app.FindFirstRecordByFilter("public_infos", "id != ''")
	if err != nil {
		return err
	}
	lastEmoteSync := publicSettings.GetDateTime("last_emote_sync")
	if lastEmoteSync.IsZero() || lastEmoteSync.Time().Add(1*time.Hour).Before(time.Now()) {
		go cronjobs.UpdateEmotes(app)
	}
	lastVodSync := publicSettings.GetDateTime("last_vod_sync")
	if lastVodSync.IsZero() || lastVodSync.Time().Add(1*time.Hour).Before(time.Now()) {
		go cronjobs.RunTwitchDownloads(app)
	}

	// schedule cronjobs
	scheduler := app.Cron()
	scheduler.SetTimezone(time.Now().Local().Location())
	scheduler.MustAdd("set_stream_status", "*/1 * * * *", func() {
		cronjobs.SetStreamStatus(app)
	})
	scheduler.MustAdd("update_emotes", "@hourly", func() {
		cronjobs.UpdateEmotes(app)
	})
	scheduler.MustAdd("twitch_downloads", "@hourly", func() {
		cronjobs.RunTwitchDownloads(app)
	})

	if err := deleteEmptyGameRecords(app); err != nil {
		return err
	}

	return nil
}

func deleteEmptyGameRecords(app *pocketbase.PocketBase) error {
	// delete empty games
	allEmptyGames, err := app.FindAllRecords("game", dbx.HashExp{
		"name": "",
	})
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	for _, game := range allEmptyGames {
		if er := app.Delete(game); er != nil {
			logger.Error.Println(err)
			return err
		}
	}
	return nil
}
