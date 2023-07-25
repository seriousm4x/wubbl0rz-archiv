package cronjobs

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/database"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/external_apis"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/filesystem"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/logger"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/models"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/queries"
)

func createSegmentsfromURL(input_url string, segmentsPath string, filename string, video_type string) error {
	var cmd *exec.Cmd
	if video_type == "vod" {
		cmd = exec.Command("ffmpeg",
			"-hide_banner", "-loglevel", "error", "-stats",
			"-i", input_url, "-map", "p:0", "-c", "copy",
			"-hls_playlist_type", "vod", "-hls_time", "10", "-hls_segment_filename",
			filepath.Join(segmentsPath, filename+"_%04d.ts"),
			filepath.Join(segmentsPath, filename+".m3u8"),
		)
	} else {
		cmd = exec.Command("ffmpeg",
			"-hide_banner", "-loglevel", "error", "-stats",
			"-i", input_url, "-c", "copy",
			"-hls_playlist_type", "vod", "-hls_time", "10", "-hls_segment_filename",
			filepath.Join(segmentsPath, filename+"_%04d.ts"),
			filepath.Join(segmentsPath, filename+".m3u8"),
		)
	}

	if err := cmd.Run(); err != nil {
		logger.Error.Println(cmd.Args)
		return err
	}

	return nil
}

func DownloadVods() (int, error) {
	logger.Debug.Println("[cronjob] vod download started")
	vods_downloaded := 0

	var vods []external_apis.TwitchHelixVideo
	if err := external_apis.TwitchGetHelixVideos(&vods); err != nil {
		return vods_downloaded, err
	}

	var streams external_apis.TwitchStreamResponse
	if err := external_apis.TwitchGetHelixStreams(&streams); err != nil {
		return vods_downloaded, err
	}
	if len(streams.Data) > 0 {
		return vods_downloaded, errors.New("stream is live. skipping vod download")
	}

	for _, vod := range vods {
		// skip vod if created less then 24h ago (only relevant for affiliates)
		// if !vod.CreatedAt.Before(time.Now().Add(time.Duration(-24) * time.Hour)) {
		// 	continue
		// }

		var newVod models.Vod
		newVod.Title = vod.Title
		newVod.Date = &vod.CreatedAt
		newVod.Viewcount = vod.ViewCount
		newVod.Filename = "v" + vod.ID
		newVod.Publish = true

		// check if vod already in db, if yes update viewcount and continue
		var v models.Vod
		if err := queries.GetVodByFilename(&v, newVod.Filename); err == nil {
			if e := queries.PatchVod(map[string]interface{}{"Viewcount": newVod.Viewcount}, v.UUID); e != nil {
				return vods_downloaded, err
			}
			continue
		}

		// create destination path
		vodsPath := filepath.Join("/var/www/media", "vods")
		segmentsPath := filepath.Join(vodsPath, newVod.Filename+"-segments")
		if err := os.MkdirAll(segmentsPath, 0755); err != nil && !os.IsExist(err) {
			return vods_downloaded, err
		}

		// get m3u8 playlist from twitch
		m3u8Url, err := external_apis.BuildDownloadURL(vod.ID, true)
		if err != nil {
			os.RemoveAll(segmentsPath)
			return vods_downloaded, err
		}

		// pass the m3u8 to ffmpeg to create .ts segments
		if err := createSegmentsfromURL(m3u8Url, segmentsPath, newVod.Filename, "vod"); err != nil {
			os.RemoveAll(segmentsPath)
			return vods_downloaded, err
		}

		// get metadata from m3u8
		var m filesystem.Meta
		m.Filename = newVod.Filename
		if err := filesystem.GetMetadata(segmentsPath, &m); err != nil {
			os.RemoveAll(segmentsPath)
			return vods_downloaded, err
		}
		newVod.Duration = m.Duration
		newVod.Resolution = m.Resolution
		newVod.Fps = m.Fps
		newVod.Size = m.Size

		// create thumbnails ...
		if err := filesystem.CreateThumbnails(vodsPath, newVod.Filename, newVod.Duration); err != nil {
			return vods_downloaded, err
		}

		// create vod in database
		if err := queries.AddNewVod(&newVod); err != nil {
			return vods_downloaded, err
		}

		vods_downloaded += 1
	}

	if vods_downloaded == 0 {
		var settings models.Settings
		settings.DateVodsUpdate = time.Now()
		if err := queries.PartiallyUpdateSettings(&settings); err != nil {
			return vods_downloaded, err
		}
	}

	logger.Debug.Printf("[cronjob] vods downloaded: %d", vods_downloaded)
	return vods_downloaded, nil
}

func DownloadClips() (int, error) {
	logger.Debug.Println("[cronjob] clip download started")
	clips_downloaded := 0
	clips_updated := 0

	var clips []external_apis.TwitchHelixClip
	if err := external_apis.TwitchGetHelixClips(&clips); err != nil {
		logger.Error.Println(err)
		return clips_downloaded, err
	}

	for _, clip := range clips {
		// skip clip if created less then 24h ago (only relevant for affiliates)
		// see: https://www.nbcnews.com/tech/twitch-partners-multiple-platforms-youtube-facebook-rcna44477
		// and: https://esports.gg/news/streamers/twitch-exclusivity-removal/
		// if !clip.CreatedAt.Before(time.Now().Add(time.Duration(-24) * time.Hour)) {
		//	 continue
		// }

		if clip.ViewCount < 3 {
			continue
		}

		// update game
		var game models.Game
		if err := queries.GetOneGame(&game, clip.GameID); err != nil {
			logger.Error.Println(err)
			game.UUID = clip.GameID
			if e := queries.AddNewGame(&game); e != nil {
				logger.Error.Println(e)
				return clips_downloaded, e
			}
		}

		// get or create creator
		var creator models.Creator
		if err := queries.GetOneCreator(&creator, clip.CreatorID); err != nil {
			// creator doesn't exist, create it
			creator.UUID = clip.CreatorID
			creator.Name = clip.CreatorName
			creator.Clips = nil
			if e := queries.AddNewCreator(&creator); e != nil {
				logger.Error.Println(e)
				return clips_downloaded, e
			}
		}

		// check if clip already in db
		var c models.Clip
		if err := queries.GetClipByFilename(&c, clip.ID); err == nil {
			// update clip title and viewcount
			changes := map[string]interface{}{
				"Title":     clip.Title,
				"Viewcount": clip.ViewCount,
			}
			if e := queries.PatchClip(changes, c.UUID); e != nil {
				logger.Error.Println(e)
				return clips_downloaded, e
			}
			clips_updated += 1
			continue
		}

		// define new clip
		var newClip models.Clip
		if clip.VideoID != "" {
			var v models.Vod
			if err := queries.GetVodByFilename(&v, "v"+clip.VideoID); err == nil {
				newClip.VodUUID = v.UUID
			}
		}
		newClip.Title = clip.Title
		newClip.Date = &clip.CreatedAt
		newClip.Filename = clip.ID
		newClip.Viewcount = clip.ViewCount
		newClip.GameUUID = game.UUID
		newClip.CreatorUUID = creator.UUID
		newClip.VodOffset = clip.VodOffset

		// create destination path
		clipsPath := filepath.Join("/var/www/media", "clips")
		segmentsPath := filepath.Join(clipsPath, newClip.Filename+"-segments")
		if err := os.MkdirAll(segmentsPath, 0755); err != nil && !os.IsExist(err) {
			logger.Error.Println(err)
			return clips_downloaded, err
		}

		// get clip url from twitch
		downloadURL, err := external_apis.BuildDownloadURL(clip.ID, false)
		if err != nil {
			logger.Error.Println(err)
			os.RemoveAll(segmentsPath)
			return clips_downloaded, err
		}

		// pass the clip url to ffmpeg to create .ts segments
		if err := createSegmentsfromURL(downloadURL, segmentsPath, clip.ID, "clip"); err != nil {
			logger.Error.Println(err)
			os.RemoveAll(segmentsPath)
			return clips_downloaded, err
		}

		// get metadata from clip url
		var m filesystem.Meta
		m.Filename = clip.ID
		if err := filesystem.GetMetadata(segmentsPath, &m); err != nil {
			logger.Error.Println(err)
			os.RemoveAll(segmentsPath)
			return clips_downloaded, err
		}
		newClip.Duration = m.Duration
		newClip.Resolution = m.Resolution
		newClip.Size = m.Size
		newClip.Fps = m.Fps

		// create thumbnails ...
		if err := filesystem.CreateThumbnails(clipsPath, clip.ID, newClip.Duration); err != nil {
			logger.Error.Println(err)
			return clips_downloaded, err
		}

		// create clip in database
		if err := queries.AddNewClip(&newClip); err != nil {
			logger.Error.Println(err)
			return clips_downloaded, err
		}

		clips_downloaded += 1
	}

	logger.Debug.Printf("[cronjob] clips downloaded: %d", clips_downloaded)
	logger.Debug.Printf("[cronjob] clips updated: %d", clips_updated)
	return clips_downloaded, nil
}

func DownloadGames() error {
	var games []models.Game

	// get all games from db
	if result := database.DB.Find(&games); result.Error != nil {
		return result.Error
	}

	// get game name and art box from twich
	requestedGames, err := external_apis.TwitchGetHelixGames(games)
	if err != nil {
		return err
	}

	gamesPath := filepath.Join("/var/www/media", "games")
	if err := os.MkdirAll(gamesPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	// save to db
	for _, game := range requestedGames {
		box_art_url := strings.Replace(game.BoxArtURL, "{width}", "100", 1)
		box_art_url = strings.Replace(box_art_url, "{height}", "133", 1)

		// download image
		out, err := os.Create(filepath.Join(gamesPath, game.ID+".jpg"))
		if err != nil {
			return err
		}
		defer out.Close()

		resp, err := http.Get(box_art_url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("box art download: status code was %d", resp.StatusCode)
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}

		if err := queries.PatchGame(map[string]interface{}{"Name": game.Name, "BoxartURL": box_art_url}, game.ID); err != nil {
			return err
		}
	}
	return nil
}
