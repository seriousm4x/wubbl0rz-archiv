package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// The twitch response for a stream request.
type TwitchStreamResponse struct {
	Data []struct {
		UserName     string    `json:"user_name"`
		UserLogin    string    `json:"user_login"`
		GameName     string    `json:"game_name"`
		GameID       string    `json:"game_id"`
		Title        string    `json:"title"`
		ViewerCount  uint      `json:"viewer_count"`
		StartedAt    time.Time `json:"started_at"`
		ThumbnailURL string    `json:"thumbnail_url"`
	}
}

// Request helix stream from twitch.
func TwitchGetHelixStreams(app *pocketbase.PocketBase, streams *TwitchStreamResponse) error {
	settings, err := app.FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	if err := TwitchUpdateBearer(app); err != nil {
		return err
	}

	url := "https://api.twitch.tv/helix/streams?user_id=" + settings.GetString("broadcaster_id")
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

	if err := json.NewDecoder(resp.Body).Decode(&streams); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
