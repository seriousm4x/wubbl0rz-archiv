package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

type TwitchBadgeVersion struct {
	ID         string `json:"id"`
	ImageURL1x string `json:"image_url_1x"`
	Title      string `json:"title"`
}

type twitchBadgeSet struct {
	SetID    string               `json:"set_id"`
	Versions []TwitchBadgeVersion `json:"versions"`
}

type twitchBadgeResponse struct {
	Data []twitchBadgeSet `json:"data"`
}

type cachedTwitchBadges struct {
	badges    map[string]TwitchBadgeVersion
	expiresAt time.Time
}

var twitchBadgeCache = struct {
	sync.Mutex
	byBroadcasterID map[string]cachedTwitchBadges
}{
	byBroadcasterID: make(map[string]cachedTwitchBadges),
}

const twitchBadgeCacheDuration = 24 * time.Hour

// TwitchGetChatBadges returns global and broadcaster-specific badges indexed by set/version.
func TwitchGetChatBadges(app core.App, broadcasterID string) (map[string]TwitchBadgeVersion, error) {
	twitchBadgeCache.Lock()
	cached, ok := twitchBadgeCache.byBroadcasterID[broadcasterID]
	twitchBadgeCache.Unlock()
	if ok && cached.expiresAt.After(time.Now()) {
		return cached.badges, nil
	}

	if err := TwitchUpdateBearer(app); err != nil {
		return nil, err
	}

	settings, err := app.FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	global, err := twitchGetBadgeSets(settings, "https://api.twitch.tv/helix/chat/badges/global")
	if err != nil {
		return nil, err
	}
	channel, err := twitchGetBadgeSets(settings, "https://api.twitch.tv/helix/chat/badges?broadcaster_id="+broadcasterID)
	if err != nil {
		return nil, err
	}

	badges := make(map[string]TwitchBadgeVersion)
	for _, badgeSet := range append(global, channel...) {
		for _, version := range badgeSet.Versions {
			badges[badgeSet.SetID+"/"+version.ID] = version
		}
	}

	twitchBadgeCache.Lock()
	twitchBadgeCache.byBroadcasterID[broadcasterID] = cachedTwitchBadges{
		badges:    badges,
		expiresAt: time.Now().Add(twitchBadgeCacheDuration),
	}
	twitchBadgeCache.Unlock()

	return badges, nil
}

func twitchGetBadgeSets(settings *core.Record, url string) ([]twitchBadgeSet, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-ID", settings.GetString("ttv_client_id"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", settings.GetString("ttv_bearer_token")))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Twitch badge request returned status %d", resp.StatusCode)
		logger.Error.Println(err)
		return nil, err
	}

	var decoded twitchBadgeResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return decoded.Data, nil
}
