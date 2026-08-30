package routes

import (
	"net/http"
	"regexp"

	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/external"
)

var twitchBroadcasterIDPattern = regexp.MustCompile(`^[0-9]+$`)

// TwitchBadges proxies Twitch chat badge metadata without exposing app credentials.
func TwitchBadges(app core.App, e *core.RequestEvent) error {
	broadcasterID := e.Request.URL.Query().Get("broadcaster_id")
	if !twitchBroadcasterIDPattern.MatchString(broadcasterID) {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"message": "broadcaster_id must be a numeric Twitch user ID",
		})
	}

	badges, err := external.TwitchGetChatBadges(app, broadcasterID)
	if err != nil {
		return e.JSON(http.StatusBadGateway, map[string]string{
			"message": "unable to load Twitch badges",
		})
	}

	return e.JSON(http.StatusOK, badges)
}
