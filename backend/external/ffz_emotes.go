package external

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

type ffzEmoteResponse struct {
	Room struct {
		Set int
	}
	Sets map[int]struct {
		Emoticons []struct {
			ID   int
			Name string
		}
	}
}

func FfzUpdateEmotes(app *pocketbase.PocketBase) error {
	settings, err := app.FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	url := "https://api.frankerfacez.com/v1/room/id/" + settings.GetString("broadcaster_id")
	resp, err := http.Get(url)
	if err != nil {
		logger.Error.Println(err)
		logger.Error.Printf("%+v", resp)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code was %d", resp.StatusCode)
		logger.Error.Printf(err.Error())
		logger.Error.Printf("%+v", resp)
		return err
	}

	var responseJson ffzEmoteResponse
	json.NewDecoder(resp.Body).Decode(&responseJson)

	collection, err := app.FindCollectionByNameOrId("emote")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	// iterate emotes and save in db
	for _, respEmote := range responseJson.Sets[responseJson.Room.Set].Emoticons {
		emote, err := app.FindFirstRecordByFilter("emote",
			"name={:name} && provider='ffz'",
			dbx.Params{"name": respEmote.Name},
		)
		if err == sql.ErrNoRows {
			emote_id := strconv.FormatInt(int64(respEmote.ID), 10)
			emote = core.NewRecord(collection)
			emote.Set("name", respEmote.Name)
			emote.Set("url", fmt.Sprintf("https://cdn.frankerfacez.com/emote/%s/4", emote_id))
			emote.Set("provider", "ffz")
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
