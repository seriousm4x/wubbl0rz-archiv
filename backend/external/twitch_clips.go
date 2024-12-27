package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// A twitch helix clip
type TwitchHelixClip struct {
	ID          string
	CreatorID   string `json:"creator_id"`
	CreatorName string `json:"creator_name"`
	GameID      string `json:"game_id"`
	Title       string
	ViewCount   int       `json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	VideoID     string    `json:"video_id"`
	VodOffset   int       `json:"vod_offset"`
}

// The twitch response for a clip request.
type TwitchHelixClipResponse struct {
	Data       []TwitchHelixClip
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// Easily calculate the date x weeks back from given time
func weeksBack(date time.Time, weeks int) time.Time {
	return date.AddDate(0, 0, -7*weeks)
}

// Request helix clips from twitch.
func TwitchGetHelixClips(app *pocketbase.PocketBase, clips *[]TwitchHelixClip) error {
	settings, err := app.FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	if err := TwitchUpdateBearer(app); err != nil {
		return err
	}

	// twitch api only returns ~ 1000-1100 clips via pagination cursor.
	// the only way to get all clips is to specify a time range (e.g. 1 week) and iterate with pagination cursor
	// through that week, then specify the next time range, iterate through that week, and so on until channel creation
	// date is reached. jeff hates this trick.
	var helixUser TwitchHelixUserResponse
	var clipResponse TwitchHelixClipResponse

	if err := TwitchGetHelixUser(app, &helixUser); err != nil {
		return err
	}

	now := time.Now()
	weeks := 1
	url := fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%s&first=100&started_at=%s",
		settings.GetString("broadcaster_id"), weeksBack(now, weeks).UTC().Format(time.RFC3339))

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
			logger.Error.Printf(err.Error())
			logger.Error.Printf("%+v", resp)
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&clipResponse); err != nil {
			logger.Error.Println(err)
			return err
		}
		*clips = append(*clips, clipResponse.Data...)

		if clipResponse.Pagination.Cursor != "" {
			url = fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%s&first=100&started_at=%s&after=%s",
				settings.GetString("broadcaster_id"), weeksBack(now, weeks).UTC().Format(time.RFC3339), clipResponse.Pagination.Cursor)
		} else {
			weeks += 1
			newWeek := weeksBack(now, weeks)
			if newWeek.Before(helixUser.Data[0].CreatedAt) {
				break
			}
			url = fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%s&first=100&started_at=%s",
				settings.GetString("broadcaster_id"), weeksBack(now, weeks).UTC().Format(time.RFC3339))
		}
	}

	return nil
}
