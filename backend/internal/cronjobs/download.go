package cronjobs

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/external"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/assets"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

var twitchDownloadsRunning = false

// Download vods, clips and games from twitch
func RunTwitchDownloads(app *pocketbase.PocketBase) {
	if twitchDownloadsRunning {
		return
	}

	twitchDownloadsRunning = true
	downloaded_items := 0

	downloaded_items += DownloadVods(app)
	downloaded_items += DownloadClips(app)
	DownloadGames(app)

	logger.Info.Printf("[cronjob] Downloaded %d items", downloaded_items)

	twitchDownloadsRunning = false
}

// Create .ts segments and .m3u8 playlist
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
		logger.Error.Println(err)
		logger.Error.Println(cmd.Args)
		return err
	}

	return nil
}

// Download all vods
func DownloadVods(app *pocketbase.PocketBase) int {
	logger.Debug.Println("[cronjob] vod download started")
	vods_downloaded := 0

	var streams external.TwitchStreamResponse
	if err := external.TwitchGetHelixStreams(app, &streams); err != nil {
		return vods_downloaded
	}

	if len(streams.Data) > 0 {
		logger.Warning.Println("Stream is live. Skipping vod download")
		return vods_downloaded
	}

	collection, err := app.FindCollectionByNameOrId("vod")
	if err != nil {
		logger.Error.Println(err)
		return vods_downloaded
	}

	var vods []external.TwitchHelixVideo
	if err := external.TwitchGetHelixVideos(app, &vods); err != nil {
		return vods_downloaded
	}

	vodsPath := path.Join(assets.ArchiveDir, "vods")

	for _, vod := range vods {
		// skip vod if created less then 24h ago (only relevant for affiliates)
		// if !vod.CreatedAt.Before(time.Now().Add(time.Duration(-24) * time.Hour)) {
		// 	continue
		// }

		var m assets.Meta
		m.Filename = "v" + vod.ID

		// check if vod already in db and update
		record, err := app.FindFirstRecordByData("vod", "filename", m.Filename)
		if err == nil {
			record.Set("title", vod.Title)
			record.Set("viewcount", vod.ViewCount)
			if err := app.Save(record); err != nil {
				logger.Error.Println(err)
				return vods_downloaded
			}
			continue
		} else if err != sql.ErrNoRows {
			logger.Error.Println(err)
			return vods_downloaded
		}

		// create new vod
		newVod := core.NewRecord(collection)
		newVod.Set("title", vod.Title)
		newVod.Set("date", vod.CreatedAt)
		newVod.Set("viewcount", vod.ViewCount)
		newVod.Set("filename", m.Filename)
		newVod.Set("publish", true)

		// create destination path
		segmentsPath := filepath.Join(vodsPath, m.Filename+"-segments")
		if err := os.MkdirAll(segmentsPath, 0755); err != nil && !os.IsExist(err) {
			logger.Error.Println(err)
			return vods_downloaded
		}
		time.Sleep(10 * time.Second)

		// get m3u8 playlist from twitch
		m3u8Url, err := external.BuildDownloadURL(vod.ID, true)
		if err != nil {
			if err := os.RemoveAll(segmentsPath); err != nil {
				logger.Error.Println(err)
			}
			continue
		}

		// pass the m3u8 to ffmpeg to create .ts segments
		if err := createSegmentsfromURL(m3u8Url, segmentsPath, m.Filename, "vod"); err != nil {
			if err := os.RemoveAll(segmentsPath); err != nil {
				logger.Error.Println(err)
			}
			continue
		}

		// get metadata from m3u8
		if err := assets.GetMetadata(segmentsPath, &m); err != nil {
			if err := os.RemoveAll(segmentsPath); err != nil {
				logger.Error.Println(err)
			}
			continue
		}
		newVod.Set("duration", m.Duration)
		newVod.Set("resolution", m.Resolution)
		newVod.Set("fps", m.Fps)
		newVod.Set("size", m.Size)

		// create vod in database
		if err := app.Save(newVod); err != nil {
			logger.Error.Println(err)
			if err := os.RemoveAll(segmentsPath); err != nil {
				logger.Error.Println(err)
			}
			return vods_downloaded
		}

		// create thumbnails
		if err := assets.CreatePreviewThumbnailsSprites(app, []string{newVod.Id}); err != nil {
			if err := os.RemoveAll(segmentsPath); err != nil {
				logger.Error.Println(err)
			}
			continue
		}

		vods_downloaded += 1
	}

	// set timestamp
	publicInfos, err := app.FindFirstRecordByFilter("public_infos", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return vods_downloaded
	}
	publicInfos.Set("last_vod_sync", time.Now())
	if err := app.Save(publicInfos); err != nil {
		logger.Error.Println(err)
		return vods_downloaded
	}

	logger.Debug.Printf("[cronjob] vods downloaded: %d", vods_downloaded)
	return vods_downloaded
}

// Download all clips
func DownloadClips(app *pocketbase.PocketBase) int {
	logger.Debug.Println("[cronjob] clip download started")
	clips_downloaded := 0

	var clips []external.TwitchHelixClip
	if err := external.TwitchGetHelixClips(app, &clips); err != nil {
		logger.Error.Println(err)
		return clips_downloaded
	}

	clipCollection, err := app.FindCollectionByNameOrId("clip")
	if err != nil {
		logger.Error.Println(err)
		return clips_downloaded
	}
	gameCollection, err := app.FindCollectionByNameOrId("game")
	if err != nil {
		logger.Error.Println(err)
		return clips_downloaded
	}
	creatorCollection, err := app.FindCollectionByNameOrId("creator")
	if err != nil {
		logger.Error.Println(err)
		return clips_downloaded
	}

	clipsPath := path.Join(assets.ArchiveDir, "clips")

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

		var m assets.Meta
		m.Filename = clip.ID

		// update game
		game, err := app.FindFirstRecordByData("game", "ttv_id", clip.GameID)
		if err == sql.ErrNoRows {
			game = core.NewRecord(gameCollection)
			game.Set("ttv_id", clip.GameID)
			game.Set("name", "Unknown")
			game.Set("box_art_url", "https://")
			if err := app.Save(game); err != nil {
				logger.Error.Println(err)
				return clips_downloaded
			}
		} else if err != nil {
			logger.Error.Println(err)
			return clips_downloaded
		}

		// get or create creator
		creator, err := app.FindFirstRecordByData("creator", "ttv_id", clip.CreatorID)
		if err == sql.ErrNoRows {
			creator = core.NewRecord(creatorCollection)
			creator.Set("ttv_id", clip.CreatorID)
			creator.Set("name", clip.CreatorName)
			if err := app.Save(creator); err != nil {
				logger.Error.Println(err)
				return clips_downloaded
			}
		} else if err != nil {
			logger.Error.Println(err)
			return clips_downloaded
		}

		// check if clip already in db and update
		record, err := app.FindFirstRecordByData("clip", "filename", clip.ID)
		if err == nil {
			record.Set("title", clip.Title)
			record.Set("viewcount", clip.ViewCount)
			if err := app.Save(record); err != nil {
				logger.Error.Println(err)
				return clips_downloaded
			}
			continue
		}

		// define new clip
		newClip := core.NewRecord(clipCollection)
		if clip.VideoID != "" {
			relatedVod, err := app.FindFirstRecordByData("vod", "filename", "v"+clip.VideoID)
			if err == nil {
				newClip.Set("vod", relatedVod.Id)
			} else if err != sql.ErrNoRows {
				logger.Error.Println(err)
				return clips_downloaded
			}
		}
		newClip.Set("title", clip.Title)
		newClip.Set("date", clip.CreatedAt)
		newClip.Set("filename", clip.ID)
		newClip.Set("viewcount", clip.ViewCount)
		newClip.Set("vod_offset", clip.VodOffset)
		newClip.Set("game", game.Id)
		newClip.Set("creator", creator.Id)

		// create destination path
		segmentsPath := filepath.Join(clipsPath, clip.ID+"-segments")
		if err := os.MkdirAll(segmentsPath, 0755); err != nil && !os.IsExist(err) {
			logger.Error.Println(err)
			return clips_downloaded
		}

		// get clip url from twitch
		downloadURL, err := external.BuildDownloadURL(clip.ID, false)
		if err != nil {
			if err := os.RemoveAll(segmentsPath); err != nil {
				logger.Error.Println(err)
			}
			continue
		}

		// pass the clip url to ffmpeg to create .ts segments
		if err := createSegmentsfromURL(downloadURL, segmentsPath, clip.ID, "clip"); err != nil {
			if err := os.RemoveAll(segmentsPath); err != nil {
				logger.Error.Println(err)
			}
			continue
		}

		// get metadata from clip url
		if err := assets.GetMetadata(segmentsPath, &m); err != nil {
			if err := os.RemoveAll(segmentsPath); err != nil {
				logger.Error.Println(err)
			}
			continue
		}
		newClip.Set("duration", m.Duration)
		newClip.Set("resolution", m.Resolution)
		newClip.Set("fps", m.Fps)
		newClip.Set("size", m.Size)

		// create clip in database
		if err := app.Save(newClip); err != nil {
			logger.Error.Println(err)
			return clips_downloaded
		}

		// create thumbnails
		if err := assets.CreatePreviewThumbnailsSprites(app, []string{newClip.Id}); err != nil {
			if err := os.RemoveAll(segmentsPath); err != nil {
				logger.Error.Println(err)
			}
			return clips_downloaded
		}

		clips_downloaded += 1
	}

	logger.Debug.Printf("[cronjob] clips downloaded: %d", clips_downloaded)
	return clips_downloaded
}

// Download all games
func DownloadGames(app *pocketbase.PocketBase) {
	logger.Debug.Println("[cronjob] game download started")

	// get all games from db
	games, err := app.FindAllRecords("game")
	if err != nil {
		logger.Error.Println(err)
		return
	}

	// get game name and art box from twich
	requestedGames, err := external.TwitchGetHelixGames(app, games)
	if err != nil {
		logger.Error.Println(err)
		return
	}

	gamesPath := path.Join(assets.ArchiveDir, "games")

	if err := os.MkdirAll(gamesPath, 0755); err != nil && !os.IsExist(err) {
		logger.Error.Println(err)
		return
	}

	// save to db
	for _, game := range requestedGames {
		box_art_url := strings.Replace(game.BoxArtURL, "{width}", "100", 1)
		box_art_url = strings.Replace(box_art_url, "{height}", "133", 1)

		// download image
		out, err := os.Create(filepath.Join(gamesPath, game.ID+".jpg"))
		if err != nil {
			logger.Error.Println(err)
			return
		}
		defer out.Close()

		resp, err := http.Get(box_art_url)
		if err != nil {
			logger.Error.Println(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("status code was %d", resp.StatusCode)
			logger.Error.Println(err)
			logger.Error.Printf("%+v", resp)
			return
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			logger.Error.Println(err)
			return
		}

		record, err := app.FindFirstRecordByData("game", "ttv_id", game.ID)
		if err != nil {
			logger.Error.Println(err)
			return
		}

		record.Set("name", game.Name)
		record.Set("box_art_url", box_art_url)

		if err := app.Save(record); err != nil {
			logger.Error.Println(err)
			continue
		}
	}
}
