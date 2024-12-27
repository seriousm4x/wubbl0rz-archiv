package external

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// A twitch helix game
type TwitchHelixGame struct {
	ID        string
	Name      string
	BoxArtURL string `json:"box_art_url"`
}

// The twitch response for a game request.
type TwitchHelixGameResponse struct {
	Data []TwitchHelixGame
}

// Request helix games from twitch.
func TwitchGetHelixGames(app *pocketbase.PocketBase, games []*core.Record) ([]TwitchHelixGame, error) {
	settings, err := app.FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	if err := TwitchUpdateBearer(app); err != nil {
		return nil, err
	}

	// twitch allows 100 games per request. so we split them in slices of 100 each and request them
	var divided [][]*core.Record
	chunkSize := 100
	for i := 0; i < len(games); i += chunkSize {
		end := i + chunkSize
		if end > len(games) {
			end = len(games)
		}
		divided = append(divided, (games)[i:end])
	}

	var gameResponse TwitchHelixGameResponse
	var respondedGames []TwitchHelixGame

	for _, chunk := range divided {
		url := "https://api.twitch.tv/helix/games?"
		for i, game := range chunk {
			if i == 0 {
				url += "id=" + game.GetString("ttv_id")
			} else {
				url += "&id=" + game.GetString("ttv_id")
			}
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Error.Println(err)
			return nil, err
		}

		req.Header.Set("Client-ID", settings.GetString("ttv_client_id"))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.GetString("ttv_bearer_token")))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error.Println(err)
			logger.Error.Printf("%+v", resp)
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("status code was %d", resp.StatusCode)
			logger.Error.Println(err)
			logger.Error.Printf("%+v", resp)
			return nil, err
		}

		if err := json.NewDecoder(resp.Body).Decode(&gameResponse); err != nil {
			logger.Error.Println(err)
			return nil, err
		}

		respondedGames = append(respondedGames, gameResponse.Data...)
	}

	return respondedGames, nil
}
