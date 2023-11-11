package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/assets"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func setVodState(app *pocketbase.PocketBase, vod *models.Record, state string) error {
	vod.Set("youtube_upload", state)
	if err := app.Dao().SaveRecord(vod); err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

func getClient(app *pocketbase.PocketBase, scope string) (*http.Client, error) {
	ctx := context.Background()

	settings, err := app.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	clientSecret := settings.GetString("yt_client_secret")
	if clientSecret == "" {
		err := fmt.Errorf("yt_client_secret is empty")
		logger.Error.Println(err)
		return nil, err
	}

	config, err := google.ConfigFromJSON([]byte(clientSecret), scope)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

	tok := &oauth2.Token{}

	bearerToken := settings.GetString("yt_bearer_token")
	if bearerToken == "" {
		err := fmt.Errorf("yt_bearer_token is empty")
		logger.Error.Println(err)
		return nil, err
	}

	if err := json.Unmarshal([]byte(bearerToken), &tok); err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return config.Client(ctx, tok), nil
}

func YoutubeUpload(app *pocketbase.PocketBase, id string) error {
	vod, err := app.Dao().FindRecordById("vod", id)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	if err := setVodState(app, vod, "pending"); err != nil {
		return err
	}

	client, err := getClient(app, youtube.YoutubeUploadScope)
	if err != nil {
		setVodState(app, vod, "")
		return err
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		logger.Error.Println(err)
		setVodState(app, vod, "")
		return err
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:      vod.GetString("title"),
			CategoryId: "28", // Science & Technology, https://gist.github.com/dgp/1b24bf2961521bd75d6c
		},
		Status: &youtube.VideoStatus{PrivacyStatus: "unlisted"},
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	// create temp video file
	filename := vod.GetString("filename")
	tempDir := filepath.Join(assets.ArchiveDir, "_temp")
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.Mkdir(tempDir, 700); err != nil {
			logger.Error.Println(err)
			setVodState(app, vod, "")
			return err
		}
	}

	tempVideo := filepath.Join(tempDir, id+".mp4")

	cmd := exec.Command("ffmpeg",
		"-i", filepath.Join(assets.ArchiveDir, "vods", filename+"-segments", filename+".m3u8"),
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		"-movflags", "frag_keyframe+empty_moov",
		"-f", "mp4", "-y", tempVideo)
	if err := cmd.Run(); err != nil {
		logger.Error.Println(err)
		setVodState(app, vod, "")
		if err := os.Remove(tempVideo); err != nil {
			logger.Error.Println(err)
			return err
		}
		return err
	}

	file, err := os.Open(tempVideo)
	defer file.Close()
	if err != nil {
		logger.Error.Println(err)
		setVodState(app, vod, "")
		if err := os.Remove(tempVideo); err != nil {
			logger.Error.Println(err)
			return err
		}
		return err
	}

	youtubeVideo, err := call.Media(file).Do()
	if err != nil {
		logger.Error.Println(err)
		setVodState(app, vod, "")
		if err := os.Remove(tempVideo); err != nil {
			logger.Error.Println(err)
			return err
		}
		return err
	}
	logger.Info.Printf("Upload successful! Video ID: %v\n", youtubeVideo.Id)

	if err := os.Remove(tempVideo); err != nil {
		logger.Error.Println(err)
		setVodState(app, vod, "")
		return err
	}

	if err := setVodState(app, vod, "done"); err != nil {
		return err
	}

	thumb := vod.GetString("custom_thumbnail")
	if thumb == "" {
		return nil
	}

	collection, err := app.Dao().FindCollectionByNameOrId("vod")
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	thumbPath := path.Join(app.DataDir(), "storage", collection.Id, vod.Id, thumb)
	thumbReader, err := os.Open(thumbPath)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	defer thumbReader.Close()

	_, err = service.Thumbnails.Set(youtubeVideo.Id).Media(thumbReader).Do()
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
