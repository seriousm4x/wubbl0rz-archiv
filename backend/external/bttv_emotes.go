package external

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

type bttvEmoteResponse struct {
	SharedEmotes []struct {
		ID   string
		Code string
	}
}

func BttvUpdateEmotes(app *pocketbase.PocketBase) error {
	settings, err := app.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	url := "https://api.betterttv.net/3/cached/users/twitch/" + settings.GetString("broadcaster_id")
	resp, err := http.Get(url)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code was %d", resp.StatusCode)
		logger.Error.Printf(err.Error())
		logger.Error.Printf("%+v", resp)
		return err
	}

	var responseJson bttvEmoteResponse
	json.NewDecoder(resp.Body).Decode(&responseJson)

	collection, err := app.Dao().FindCollectionByNameOrId("emote")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	for _, respEmote := range responseJson.SharedEmotes {
		emote, err := app.Dao().FindFirstRecordByFilter("emote",
			"name={:name} && provider='bttv'",
			dbx.Params{"name": respEmote.Code},
		)
		if err == sql.ErrNoRows {
			emote = models.NewRecord(collection)
			emote.Set("name", respEmote.Code)
			emote.Set("url", fmt.Sprintf("https://cdn.betterttv.net/emote/%s/3x", respEmote.ID))
			emote.Set("provider", "bttv")
		} else if err != nil {
			logger.Error.Println(err)
			return err
		} else {
			emote.Set("outdated", false)
		}

		if err := app.Dao().SaveRecord(emote); err != nil {
			logger.Error.Println(err)
			return err
		}
	}

	return nil
}
