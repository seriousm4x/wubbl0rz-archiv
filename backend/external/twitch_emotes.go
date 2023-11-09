package external

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// The twitch response for an emote request.
type twitchEmoteResponse struct {
	Data []struct {
		ID     string
		Name   string
		Images struct {
			Url4x string `json:"url_4x"`
		}
		Format []string
	}
}

// Update all twitch emotes
func TwitchUpdateEmotes(app *pocketbase.PocketBase) error {
	settings, err := app.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	if err := TwitchUpdateBearer(app); err != nil {
		return err
	}

	url := "https://api.twitch.tv/helix/chat/emotes?broadcaster_id=" + settings.GetString("broadcaster_id")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	req.Header.Set("Client-ID", settings.GetString("ttv_client_id"))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.GetString("ttv_bearer_token")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code was %d", resp.StatusCode)
		logger.Error.Printf(err.Error())
		logger.Error.Printf("%+v", req)
		return err
	}

	var respDecoded twitchEmoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&respDecoded); err != nil {
		logger.Error.Println(err)
		return err
	}

	collection, err := app.Dao().FindCollectionByNameOrId("emote")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	for _, respEmote := range respDecoded.Data {
		var image string
		if slices.Contains(respEmote.Format, "animated") {
			image = strings.Replace(respEmote.Images.Url4x, "/static/", "/animated/", 1)
		} else {
			image = respEmote.Images.Url4x
		}

		emote, err := app.Dao().FindFirstRecordByFilter("emote",
			"name={:name} && provider='twitch'",
			dbx.Params{"name": respEmote.Name},
		)
		if err == sql.ErrNoRows {
			emote = models.NewRecord(collection)
			emote.Set("name", respEmote.Name)
			emote.Set("url", image)
			emote.Set("provider", "twitch")
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
