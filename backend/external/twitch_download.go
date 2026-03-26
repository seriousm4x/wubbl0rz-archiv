package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// Builds the download URL for vods/clips.
// This is a go implementation of streamlink's way to build the m3u8.
// https://github.com/streamlink/streamlink/blob/master/src/streamlink/plugins/twitch.py
func BuildDownloadURL(id string, isVod bool) (string, error) {
	var postBodyString string

	if isVod {
		postBodyString = fmt.Sprintf(`
		{
			"operationName": "PlaybackAccessToken",
			"extensions": {
				"persistedQuery": {
					"version": 1,
					"sha256Hash": "ed230aa1e33e07eebb8928504583da78a5173989fadfb1ac94be06a04f3cdbe9"
				}
			},
			"variables": {
				"isLive": false,
				"login": "",
				"isVod": true,
				"vodID": "%s",
				"playerType": "embed",
				"platform":   "site"
			}
		}`, id)
	} else {
		postBodyString = fmt.Sprintf(`
			{
				"operationName": "VideoAccessToken_Clip",
				"extensions": {
					"persistedQuery": {
						"version": 1,
						"sha256Hash": "993d9a5131f15a37bd16f32342c44ed1e0b1a9b968c6afdb662d2cddd595f6c5"
					}
				},
				"variables": {
					"slug": "%s",
					"platform": "web"
				}
			}`, id)
	}

	var dummy map[string]interface{}
	if err := json.Unmarshal([]byte(postBodyString), &dummy); err != nil {
		logger.Error.Println(err)
		return "", err
	}

	postBody, err := json.Marshal(dummy)
	if err != nil {
		logger.Error.Println(err)
		return "", err
	}

	req, err := http.NewRequest("POST", "https://gql.twitch.tv/gql", bytes.NewBuffer(postBody))
	if err != nil {
		logger.Error.Println(err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Client-ID", "kimne78kx3ncx6brgo4mv6wki5h1ko")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error.Println(err)
		logger.Error.Printf("%+v", resp)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code was %d", resp.StatusCode)
		logger.Error.Println(err)
		logger.Error.Printf("%+v", resp)
		return "", err
	}

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
			logger.Error.Println(err)
			return "", err
		}
		params := url.Values{}
		params.Set("nauthsig", playbackResponse.Data.VideoPlaybackAccessToken.Signature)
		params.Set("nauth", playbackResponse.Data.VideoPlaybackAccessToken.Value)
		params.Set("allow_source", "true")
		params.Set("allow_audio_only", "true")
		params.Set("playlist_include_framerate", "true")

		return fmt.Sprintf(
			"https://usher.ttvnw.net/vod/v2/%s.m3u8?%s",
			id,
			params.Encode(),
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
		logger.Error.Println(err)
		return "", err
	}

	qualities := playbackResponse.Data.Clip.VideoQualities
	sort.Slice(qualities, func(i, j int) bool {
		qi, _ := strconv.Atoi(qualities[i].Quality)
		qj, _ := strconv.Atoi(qualities[j].Quality)
		return qi > qj
	})

	best := qualities[0]

	return fmt.Sprintf("%s?token=%s&sig=%s",
		best.SourceURL,
		url.QueryEscape(playbackResponse.Data.Clip.PlaybackAccessToken.Value),
		playbackResponse.Data.Clip.PlaybackAccessToken.Signature,
	), nil
}
