package external

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

type SevenTvEmoteResponse struct {
	EmoteSet struct {
		Emotes []struct {
			ID   string
			Name string
		}
	} `json:"emote_set"`
}

func SevenTvUpdateEmotes(app *pocketbase.PocketBase) error {
	settings, err := app.FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	url := "https://7tv.io/v3/users/twitch/" + settings.GetString("broadcaster_id")
	resp, err := http.Get(url)
	if err != nil {
		logger.Error.Println(err)
		logger.Error.Printf("%+v", resp)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code was %d", resp.StatusCode)
		logger.Error.Println(err)
		logger.Error.Printf("%+v", resp)
		return err
	}

	var responseJson SevenTvEmoteResponse
	json.NewDecoder(resp.Body).Decode(&responseJson)

	collection, err := app.FindCollectionByNameOrId("emote")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	for _, respEmote := range responseJson.EmoteSet.Emotes {
		emote, err := app.FindFirstRecordByFilter("emote",
			"name={:name} && provider='7tv'",
			dbx.Params{"name": respEmote.Name},
		)
		if err == sql.ErrNoRows {
			emote = core.NewRecord(collection)
			emote.Set("name", respEmote.Name)
			emote.Set("url", fmt.Sprintf("https://cdn.7tv.app/emote/%s/4x.webp", respEmote.ID))
			emote.Set("provider", "7tv")
		} else if err != nil {
			logger.Error.Println(err)
			return err
		} else {
			emote.Set("outdated", false)
		}

		if err := app.Save(emote); err != nil {
			logger.Error.Println(err)
			return err
		}
	}

	return nil
}
