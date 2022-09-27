package cronjobs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/AgileProggers/archiv-backend-go/pkg/external_apis"
	"github.com/AgileProggers/archiv-backend-go/pkg/logger"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"
)

func SetStreamStatus() {
	// get settings from db and update bearer if needed
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		logger.Error.Println(err)
	}
	err := external_apis.UpdateBearer(&settings)
	if err != nil {
		logger.Error.Println(err)
	}

	url := "https://api.twitch.tv/helix/streams?user_id=" + settings.BroadcasterId
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error.Println(err)
	}
	req.Header.Set("Client-ID", settings.TtvClientId)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.TtvBearerToken))
	resp, err := client.Do(req)
	if err != nil {
		logger.Error.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error.Printf("[cronjob] live check: status code was %d", resp.StatusCode)
	}

	var responseJson external_apis.TwitchStreamResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseJson); err != nil {
		logger.Error.Println(err)
	}

	isLive := len(responseJson.Data) > 0

	if isLive != settings.IsLive {
		settings = models.Settings{IsLive: true}
		if err := queries.PartiallyUpdateSettings(&settings); err != nil {
			logger.Error.Println(err)
		}
		if isLive {
			logger.Debug.Println("[cronjob] stream live")
			if os.Getenv("DISCORD_WEBHOOK") != "" {
				if err := external_apis.DiscordSendWebhook(responseJson); err != nil {
					logger.Error.Println(err)
				}
			}
		} else {
			settings.IsLive = false
			if err := queries.OverwriteAllSettings(&settings); err != nil {
				logger.Error.Println(err)
			}
		}
	}
}
