package cronjobs

import (
	"database/sql"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/external"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// Sets the stream status
func SetStreamStatus(app *pocketbase.PocketBase) error {
	var streams external.TwitchStreamResponse
	if err := external.TwitchGetHelixStreams(app, &streams); err != nil {
		logger.Error.Println(err)
		return err
	}
	isLive := len(streams.Data) > 0

	publicInfos, err := app.FindFirstRecordByFilter("public_infos", "id != ''")
	if err == sql.ErrNoRows {
		collection, err := app.FindCollectionByNameOrId("public_infos")
		if err != nil {
			logger.Error.Println(err)
			return err
		}
		publicInfos = core.NewRecord(collection)
		if err := app.Save(publicInfos); err != nil {
			logger.Error.Println(err)
			return err
		}
	} else if err != nil {
		logger.Error.Println(err)
		return err
	}

	if err := external.TwitchUpdateBearer(app); err != nil {
		return err
	}

	if isLive != publicInfos.GetBool("is_live") {
		publicInfos.Set("is_live", isLive)
		if err := app.Save(publicInfos); err != nil {
			logger.Error.Println(err)
			return err
		}
		if isLive {
			logger.Debug.Println("[jobs] stream live")
			if os.Getenv("DISCORD_WEBHOOK") != "" {
				if err := external.DiscordSendWebhook(app, streams); err != nil {
					logger.Error.Println(err)
					return err
				}
			}
		}
	}

	return nil
}
