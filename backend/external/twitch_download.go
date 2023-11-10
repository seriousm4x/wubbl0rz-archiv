package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// Builds the download URL for vods/clips.
// This is a go implementation of streamlink's way build the m3u8.
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
		logger.Error.Println(err)
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
