package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// The twitch response for a helix user request.
type TwitchHelixUserResponse struct {
	Data []struct {
		ID              string
		ProfileImageUrl string    `json:"profile_image_url"`
		CreatedAt       time.Time `json:"created_at"`
	}
}

// Request helix user from twitch.
func TwitchGetHelixUser(app *pocketbase.PocketBase, user *TwitchHelixUserResponse) error {
	settings, err := app.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	if err := TwitchUpdateBearer(app); err != nil {
		return err
	}

	url := "https://api.twitch.tv/helix/users?login=" + os.Getenv("BROADCASTER_NAME")
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

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}

// Update the settings broadcaster_id field
func TwitchUpdateBroadcasterId(app *pocketbase.PocketBase) error {
	settings, err := app.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	if err := TwitchUpdateBearer(app); err != nil {
		return err
	}

	var helixUser TwitchHelixUserResponse
	if err = TwitchGetHelixUser(app, &helixUser); err != nil {
		return err
	}

	settings.Set("broadcaster_id", helixUser.Data[0].ID)

	if err := app.Dao().SaveRecord(settings); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
