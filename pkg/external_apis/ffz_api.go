package external_apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/models"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/queries"
)

type ffzEmote struct {
	ID   int
	Name string
}

type ffzEmoteResponse struct {
	Room struct {
		Set int
	}
	Sets map[int]struct {
		Emoticons []ffzEmote
	}
}

func FfzUpdateEmotes() error {
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		return err
	}

	url := "https://api.frankerfacez.com/v1/room/id/" + settings.BroadcasterId
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ffz emotes: status code was %d", resp.StatusCode)
	}

	var responseJson ffzEmoteResponse
	json.NewDecoder(resp.Body).Decode(&responseJson)

	// iterate emotes and save in db
	for _, emote := range responseJson.Sets[responseJson.Room.Set].Emoticons {
		emote_id := strconv.FormatInt(int64(emote.ID), 10)
		newEmote := models.Emote{
			ID:       emote_id,
			Name:     emote.Name,
			URL:      fmt.Sprintf("https://cdn.frankerfacez.com/emote/%s/4", emote_id),
			Provider: "ffz",
			Outdated: false,
		}
		queries.UpdateOrCreateEmote(&newEmote, emote_id)
	}

	return nil
}
