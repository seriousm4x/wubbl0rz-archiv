package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// A twitch helix video
type TwitchHelixVideo struct {
	ID        string
	Title     string
	CreatedAt time.Time `json:"created_at"`
	ViewCount int       `json:"view_count"`
}

// The twitch response for a video request.
type TwitchHelixVideoResponse struct {
	Data       []TwitchHelixVideo
	Pagination struct {
		Cursor string
	}
}

// Request helix videos from twitch.
func TwitchGetHelixVideos(app *pocketbase.PocketBase, vods *[]TwitchHelixVideo) error {
	settings, err := app.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	if err := TwitchUpdateBearer(app); err != nil {
		return err
	}

	var videoResponse TwitchHelixVideoResponse
	url := fmt.Sprintf("https://api.twitch.tv/helix/videos?user_id=%s&type=archive&first=100",
		settings.GetString("broadcaster_id"))

	for {
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

		if err := json.NewDecoder(resp.Body).Decode(&videoResponse); err != nil {
			logger.Error.Println(err)
			return err
		}

		*vods = append(*vods, videoResponse.Data...)
		if len(videoResponse.Data) < 100 {
			break
		}

		url = fmt.Sprintf("https://api.twitch.tv/helix/videos?user_id=%s&type=archive&first=100&after=%s",
			settings.GetString("broadcaster_id"), videoResponse.Pagination.Cursor)
	}

	return nil
}
