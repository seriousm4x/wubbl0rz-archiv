package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// The twitch response for a bearer request.
type twitchBearerResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint32 `json:"expires_in"`
}

// Refresh the twitch bearer token if it expires in less than 24 hours.
func TwitchUpdateBearer(app *pocketbase.PocketBase) error {
	settings, err := app.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	if !settings.GetTime("ttv_bearer_expire").IsZero() && settings.GetTime("ttv_bearer_expire").Before(time.Now().Add(24*time.Hour)) {
		return nil
	}

	body := map[string]string{
		"client_id":     settings.GetString("ttv_client_id"),
		"client_secret": settings.GetString("ttv_client_secret"),
		"grant_type":    "client_credentials",
	}

	buf, err := json.Marshal(body)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	resp, err := http.Post("https://id.twitch.tv/oauth2/token", "application/json", bytes.NewBuffer(buf))
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

	var respDecoded twitchBearerResponse
	if err := json.NewDecoder(resp.Body).Decode(&respDecoded); err != nil {
		logger.Error.Println(err)
		return err
	}

	settings.Set("ttv_bearer_token", respDecoded.AccessToken)
	settings.Set("ttv_bearer_expire", time.Now().Add(time.Duration(respDecoded.ExpiresIn)*time.Second))

	if err := app.Dao().SaveRecord(settings); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
