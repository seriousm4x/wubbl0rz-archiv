package cronjobs

import (
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/external"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// Updates all emotes in database
func UpdateEmotes(app *pocketbase.PocketBase) error {
	logger.Debug.Println("[jobs] updating all emotes")

	// set all emotes outdated
	emotes, err := app.Dao().FindRecordsByExpr("emote")
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	for _, emote := range emotes {
		emote.Set("outdated", true)
		if err := app.Dao().SaveRecord(emote); err != nil {
			logger.Error.Println(err)
			return err
		}
	}

	// fetch emotes
	if err := external.TwitchUpdateEmotes(app); err != nil {
		return err
	}
	if err := external.BttvUpdateEmotes(app); err != nil {
		return err
	}
	if err := external.FfzUpdateEmotes(app); err != nil {
		return err
	}
	if err := external.SevenTvUpdateEmotes(app); err != nil {
		return err
	}

	// delete outdated emotes
	emotes, err = app.Dao().FindRecordsByExpr("emote", dbx.HashExp{"outdated": true})
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	for _, emote := range emotes {
		if err := app.Dao().DeleteRecord(emote); err != nil {
			logger.Error.Println(err)
			return err
		}
	}

	// set timestamp
	publicInfos, err := app.Dao().FindFirstRecordByFilter("public_infos", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	publicInfos.Set("last_emote_sync", time.Now())
	if err := app.Dao().SaveRecord(publicInfos); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
