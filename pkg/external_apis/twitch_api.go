package external_apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"
)

type twitchBearerResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint32 `json:"expires_in"`
}

type twitchEmoteResponse struct {
	Data []struct {
		ID     string
		Name   string
		Images struct {
			Url4x string `json:"url_4x"`
		}
		Format []string
	}
}

type TwitchHelixUserResponse struct {
	Data []struct {
		ID              string
		ProfileImageUrl string    `json:"profile_image_url"`
		CreatedAt       time.Time `json:"created_at"`
	}
}

type TwitchHelixVideo struct {
	ID        string
	Title     string
	CreatedAt time.Time `json:"created_at"`
	ViewCount int       `json:"view_count"`
}

type TwitchHelixVideoResponse struct {
	Data       []TwitchHelixVideo
	Pagination struct {
		Cursor string
	}
}

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

type TwitchHelixClipResponse struct {
	Data       []TwitchHelixClip
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

type TwitchHelixGame struct {
	ID        string
	Name      string
	BoxArtURL string `json:"box_art_url"`
}

type TwitchHelixGameResponse struct {
	Data []TwitchHelixGame
}

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

func sliceContains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func UpdateBearer(settings *models.Settings) error {
	// check for expire date
	if settings.TtvBearerExpireDate.After(time.Now().Add(24 * time.Hour)) {
		return nil
	}

	// update if it expires in the next 24h
	body := map[string]string{
		"client_id":     settings.TtvClientId,
		"client_secret": settings.TtvClientSecret,
		"grant_type":    "client_credentials",
	}

	json_data, err := json.Marshal(body)

	if err != nil {
		return err
	}

	resp, err := http.Post("https://id.twitch.tv/oauth2/token", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ttv bearer: status code was %d", resp.StatusCode)
	}

	var responseJson twitchBearerResponse
	json.NewDecoder(resp.Body).Decode(&responseJson)

	settings.TtvBearerToken = responseJson.AccessToken
	settings.TtvBearerExpireDate = time.Now().Add(time.Duration(responseJson.ExpiresIn) * time.Second)
	if err := queries.PartiallyUpdateSettings(settings); err != nil {
		return err
	}

	return nil
}

func TwitchGetHelixUser(user *TwitchHelixUserResponse) error {
	// get settings from db and update bearer if needed
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		return err
	}
	err := UpdateBearer(&settings)
	if err != nil {
		return err
	}

	// get twitch channel user infos
	url := "https://api.twitch.tv/helix/users?login=" + os.Getenv("BROADCASTER_NAME")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Client-ID", settings.TtvClientId)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.TtvBearerToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ttv broadcaster id: status code was %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return err
	}

	return nil
}

func TwitchUpdateBroadcasterID() error {
	// get settings from db and update bearer if needed
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		return err
	}
	err := UpdateBearer(&settings)
	if err != nil {
		return err
	}

	// get user info from twitch
	var helixUser TwitchHelixUserResponse
	if err = TwitchGetHelixUser(&helixUser); err != nil {
		return err
	}

	// save id in db
	settings.BroadcasterId = helixUser.Data[0].ID
	if err := queries.OverwriteAllSettings(&settings); err != nil {
		return err
	}

	return nil
}

func TwitchUpdateEmotes() error {
	// get settings from db and update bearer if needed
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		return err
	}
	err := UpdateBearer(&settings)
	if err != nil {
		return err
	}

	// get channel emotes
	url := "https://api.twitch.tv/helix/chat/emotes?broadcaster_id=" + settings.BroadcasterId
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Client-ID", settings.TtvClientId)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.TtvBearerToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ttv emotes: status code was %d", resp.StatusCode)
	}

	var responseJson twitchEmoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseJson); err != nil {
		return err
	}

	// iterate emotes and save in db
	for _, emote := range responseJson.Data {
		var image string
		if sliceContains(emote.Format, "animated") {
			image = strings.Replace(emote.Images.Url4x, "/static/", "/animated/", 1)
		} else {
			image = emote.Images.Url4x
		}
		newEmote := models.Emote{
			ID:       emote.ID,
			Name:     emote.Name,
			URL:      image,
			Provider: "twitch",
			Outdated: false,
		}
		queries.UpdateOrCreateEmote(&newEmote, emote.ID)
	}
	return nil
}

func TwitchGetHelixVideos(vods *[]TwitchHelixVideo) error {
	// get settings from db and update bearer if needed
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		return err
	}
	err := UpdateBearer(&settings)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.twitch.tv/helix/videos?user_id=%s&type=archive&first=100", settings.BroadcasterId)
	var videoResponse TwitchHelixVideoResponse

	for {
		// get all vods from api
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Client-ID", settings.TtvClientId)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.TtvBearerToken))
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("helix videos: status code was %d", resp.StatusCode)
		}
		if err := json.NewDecoder(resp.Body).Decode(&videoResponse); err != nil {
			return err
		}
		*vods = append(*vods, videoResponse.Data...)
		if len(videoResponse.Data) < 100 {
			break
		}
		url = fmt.Sprintf("https://api.twitch.tv/helix/videos?user_id=%s&type=archive&first=100&after=%s", settings.BroadcasterId, videoResponse.Pagination.Cursor)
	}

	return nil
}

func weeksBack(date time.Time, weeks int) time.Time {
	return date.AddDate(0, 0, -7*weeks)
}

func TwitchGetHelixClips(clips *[]TwitchHelixClip) error {
	// get settings from db and update bearer if needed
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		return err
	}
	err := UpdateBearer(&settings)
	if err != nil {
		return err
	}

	// twitch api only returns ~ 1000-1100 clips via pagination cursor.
	// the only way to get all clips is to specify a time range (e.g. 1 week) and iterate with pag cursor through that week,
	// then specify the next time range, iterate through that week, and so on until channel creation date is reached.
	// jeff hates this trick
	var helixUser TwitchHelixUserResponse
	if err := TwitchGetHelixUser(&helixUser); err != nil {
		return err
	}

	now := time.Now()
	weeks := 1

	var clipResponse TwitchHelixClipResponse
	url := fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%s&first=100&started_at=%s",
		settings.BroadcasterId, weeksBack(now, weeks).UTC().Format(time.RFC3339))

	for {
		// get all clips from api
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Client-ID", settings.TtvClientId)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.TtvBearerToken))
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("helix clips: status code was %d", resp.StatusCode)
		}
		if err := json.NewDecoder(resp.Body).Decode(&clipResponse); err != nil {
			return err
		}
		*clips = append(*clips, clipResponse.Data...)

		// check for pagination cursor. if none, set time range 1 week back
		if clipResponse.Pagination.Cursor != "" {
			url = fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%s&first=100&started_at=%s&after=%s",
				settings.BroadcasterId, weeksBack(now, weeks).UTC().Format(time.RFC3339), clipResponse.Pagination.Cursor)
		} else {
			weeks += 1
			newWeek := weeksBack(now, weeks)
			if newWeek.Before(helixUser.Data[0].CreatedAt) {
				break
			}
			url = fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%s&first=100&started_at=%s",
				settings.BroadcasterId, weeksBack(now, weeks).UTC().Format(time.RFC3339))
		}
	}

	return nil
}

func TwitchGetHelixGames(games []models.Game) ([]TwitchHelixGame, error) {
	var respondedGames []TwitchHelixGame

	// get settings from db and update bearer if needed
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		return respondedGames, err
	}
	err := UpdateBearer(&settings)
	if err != nil {
		return respondedGames, err
	}

	// twitch allows 100 games per request. so we split them in slices of 100 each and request them
	var divided [][]models.Game
	chunkSize := 100
	for i := 0; i < len(games); i += chunkSize {
		end := i + chunkSize
		if end > len(games) {
			end = len(games)
		}
		divided = append(divided, (games)[i:end])
	}

	// request chunks
	var gameResponse TwitchHelixGameResponse
	for _, chunk := range divided {
		url := "https://api.twitch.tv/helix/games?"
		for i, game := range chunk {
			if i == 0 {
				url += "id=" + game.UUID
			}
			url += "&id=" + game.UUID
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return respondedGames, err
		}
		req.Header.Set("Client-ID", settings.TtvClientId)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.TtvBearerToken))
		resp, err := client.Do(req)
		if err != nil {
			return respondedGames, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return respondedGames, fmt.Errorf("helix games: status code was %d", resp.StatusCode)
		}
		if err := json.NewDecoder(resp.Body).Decode(&gameResponse); err != nil {
			return respondedGames, err
		}

		respondedGames = append(respondedGames, gameResponse.Data...)
	}

	return respondedGames, nil
}

func TwitchGetHelixStreams(streams *TwitchStreamResponse) error {
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		return err
	}

	err := UpdateBearer(&settings)
	if err != nil {
		return err
	}

	url := "https://api.twitch.tv/helix/streams?user_id=" + settings.BroadcasterId
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Client-ID", settings.TtvClientId)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.TtvBearerToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&streams); err != nil {
		return err
	}

	return nil
}

func BuildDownloadURL(id string, isVod bool) (string, error) {
	// this is a go implementation of streamlink's way to the m3u8 from a vod/clip
	// https://github.com/streamlink/streamlink/blob/master/src/streamlink/plugins/twitch.py

	var postBodyString string

	if isVod {
		postBodyString = fmt.Sprintf(`
		{
			"operationName": "PlaybackAccessToken",
			"extensions": {
				"persistedQuery": {
					"version": 1,
					"sha256Hash": "0828119ded1c13477966434e15800ff57ddacf13ba1911c129dc2200705b0712"
				}
			},
			"variables": {
				"isLive": true,
				"login": "",
				"isVod": true,
				"vodID": "%s",
				"playerType": "embed"
			}
		}`, id)
	} else {
		postBodyString = fmt.Sprintf(`
			{
				"operationName": "VideoAccessToken_Clip",
				"extensions": {
					"persistedQuery": {
						"version": 1,
						"sha256Hash": "36b89d2507fce29e5ca551df756d27c1cfe079e2609642b4390aa4c35796eb11"
					}
				},
				"variables": {
					"slug": "%s"
				}
			}`, id)
	}

	var dummy map[string]interface{}
	if err := json.Unmarshal([]byte(postBodyString), &dummy); err != nil {
		return "", err
	}
	postBody, err := json.Marshal(dummy)
	if err != nil {
		return "", err
	}

	// post to twitch graphql
	req, err := http.NewRequest("POST", "https://gql.twitch.tv/gql", bytes.NewBuffer(postBody))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Client-ID", "kimne78kx3ncx6brgo4mv6wki5h1ko") // this is streamlink's client-id (or twitch's?)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("build download url: status code was %d", resp.StatusCode)
	}

	// return vod url
	if isVod {
		var playbackResponse struct {
			Data struct {
				VideoPlaybackAccessToken struct {
					Value     string `json:"value"`
					Signature string `json:"signature"`
				} `json:"videoPlaybackAccessToken"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&playbackResponse); err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"https://usher.ttvnw.net/vod/%s.m3u8?Client-ID=kimne78kx3ncx6brgo4mv6wki5h1ko&token=%s&sig=%s&allow_source=true&allow_audio_only=true",
			id,
			playbackResponse.Data.VideoPlaybackAccessToken.Value,
			playbackResponse.Data.VideoPlaybackAccessToken.Signature,
		), nil
	}

	var playbackResponse struct {
		Data struct {
			Clip struct {
				PlaybackAccessToken struct {
					Value     string `json:"value"`
					Signature string `json:"signature"`
				} `json:"playbackAccessToken"`
				VideoQualities []struct {
					Quality   string `json:"quality"`
					SourceURL string `json:"sourceURL"`
				} `json:"videoQualities"`
			} `json:"clip"`
		} `json:"data"`
	}

	// return clip url
	if err := json.NewDecoder(resp.Body).Decode(&playbackResponse); err != nil {
		return "", err
	}

	bestQuality := struct {
		Resolution int
		URL        string
	}{
		Resolution: 0,
		URL:        "",
	}

	for _, quality := range playbackResponse.Data.Clip.VideoQualities {
		res, err := strconv.Atoi(quality.Quality)
		if err != nil {
			continue
		}
		if bestQuality.Resolution == 0 {
			bestQuality.Resolution = res
			bestQuality.URL = quality.SourceURL
			continue
		}
		if res > bestQuality.Resolution {
			bestQuality.Resolution = res
			bestQuality.URL = quality.SourceURL
			continue
		}
	}

	return fmt.Sprintf("%s?token=%s&sig=%s",
		bestQuality.URL,
		url.QueryEscape(playbackResponse.Data.Clip.PlaybackAccessToken.Value),
		playbackResponse.Data.Clip.PlaybackAccessToken.Signature,
	), nil
}
