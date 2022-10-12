package cronjobs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/database"
	"github.com/AgileProggers/archiv-backend-go/pkg/external_apis"
	"github.com/AgileProggers/archiv-backend-go/pkg/logger"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"

	"github.com/h2non/bimg"
)

type meta struct {
	Filename   string
	Duration   int
	Resolution string
	Fps        float32
	Size       int
}

func getSegmentSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.Contains(info.Name(), ".ts") {
			size += info.Size()
		}
		return err
	})
	return size, err
}

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

func getMetadata(destPath string, m *meta) error {
	// get width, height, fps and duration
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
		"program_stream=width,height,r_frame_rate:format=duration", "-of", "default=noprint_wrappers=1:nokey=1",
		filepath.Join(destPath, m.Filename+".m3u8"))

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	splittedStdout := strings.Split(stdout.String(), "\n")
	if len(splittedStdout) < 4 {
		return errors.New("not enough values to unpack")
	}
	width := strings.TrimSpace(splittedStdout[0])
	height := strings.TrimSpace(splittedStdout[1])
	fps, err := strconv.ParseFloat(strings.TrimSpace(strings.Replace(splittedStdout[2], "/1", "", 1)), 64)
	if err != nil {
		return err
	}
	duration, err := strconv.ParseFloat(strings.TrimSpace(splittedStdout[3]), 64)
	if err != nil {
		return err
	}

	if width == "" {
		return errors.New("width empty")
	} else if height == "" {
		return errors.New("height empty")
	} else if fps == 0 {
		return errors.New("fps empty")
	} else if duration == 0 {
		return errors.New("duration is 0")
	}

	m.Duration = int(math.Round(duration))
	m.Resolution = width + "x" + height
	m.Fps = float32(fps)

	// get filesize
	size, err := getSegmentSize(destPath)
	if err != nil {
		return err
	}
	m.Size = int(size)

	return nil
}

func createThumbnails(destPath string, filename string, duration int) error {
	m3u8 := filepath.Join(destPath, filename+"-segments", filename+".m3u8")

	var timecode_framegrab string
	if duration <= 10 {
		timecode_framegrab = "0"
	} else {
		timecode_framegrab = fmt.Sprintf("%d", int(duration/2))
	}

	// create lossless source png
	src_png := filepath.Join(destPath, filename+"-source.png")
	cmd := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
		"-vframes", "1", "-f", "image2", "-y", src_png)
	if err := cmd.Run(); err != nil {
		return err
	}

	// read lossless png to buff
	buffer, err := bimg.Read(src_png)
	if err != nil {
		return err
	}

	// define final thumbnails
	type Thumbnail struct {
		Options  *bimg.Options
		Filename string
	}

	thumbnails := []Thumbnail{}
	thumbnails = append(thumbnails, Thumbnail{Filename: "-sm.jpg", Options: &bimg.Options{Width: 260, Height: 146, Type: bimg.JPEG}})
	thumbnails = append(thumbnails, Thumbnail{Filename: "-md.jpg", Options: &bimg.Options{Width: 520, Height: 293, Type: bimg.JPEG}})
	thumbnails = append(thumbnails, Thumbnail{Filename: "-lg.jpg", Options: &bimg.Options{Width: 1592, Height: 896, Type: bimg.JPEG}})
	thumbnails = append(thumbnails, Thumbnail{Filename: "-sm.avif", Options: &bimg.Options{Width: 260, Height: 146, Type: bimg.AVIF}})
	thumbnails = append(thumbnails, Thumbnail{Filename: "-md.avif", Options: &bimg.Options{Width: 520, Height: 293, Type: bimg.AVIF}})

	// create defined thumbnails
	for _, thumb := range thumbnails {
		thumb.Options.Quality = 95
		newImage, err := bimg.NewImage(buffer).Process(*thumb.Options)
		if err != nil {
			return err
		}
		bimg.Write(filepath.Join(destPath, filename+thumb.Filename), newImage)
	}

	// remove source png
	if err := os.Remove(src_png); err != nil {
		return err
	}

	// animated webp
	animated_webp := filepath.Join(destPath, filename+"-preview.webp")
	cmd = exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab,
		"-i", m3u8, "-c:v", "libwebp", "-vf", "scale=260:-1,fps=fps=15", "-lossless",
		"0", "-compression_level", "3", "-q:v", "70", "-loop", "0", "-preset", "picture",
		"-an", "-vsync", "0", "-t", "4", "-y", animated_webp)
	if err := cmd.Run(); err != nil {
		return err
	}
	// check if webp is larger than 8 byte.
	// sometimes the above command fails and the image is empty.
	// we need to move the seekpoint "-ss" after the input file
	webp_info, err := os.Stat(animated_webp)
	if err != nil {
		return err
	}
	if webp_info.Size() <= 8 {
		cmd = exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-i", m3u8,
			"-ss", timecode_framegrab, "-c:v", "libwebp", "-vf", "scale=260:-1,fps=fps=15", "-lossless",
			"0", "-compression_level", "3", "-q:v", "70", "-loop", "0", "-preset", "picture",
			"-an", "-vsync", "0", "-t", "4", "-y", animated_webp)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// create sprites
	sprite_dir := filepath.Join(destPath, filename+"-sprites")
	if err := os.MkdirAll(sprite_dir, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	cmd = exec.Command("ffmpeg", "-i", m3u8, "-vf", "fps=1/20,scale=-1:90,tile",
		"-c:v", "libwebp", "-y", filepath.Join(sprite_dir, filename+"_%03d.webp"))
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func DownloadVods() error {
	logger.Debug.Println("[cronjob] vod download started")
	var vods []external_apis.TwitchHelixVideo
	if err := external_apis.TwitchGetHelixVideos(&vods); err != nil {
		logger.Error.Println(err)
	}

	vods_downloaded := 0
	for _, vod := range vods {
		// skip vod if created less then 24h ago (terms of service)
		if !vod.CreatedAt.Before(time.Now().Add(time.Duration(-24) * time.Hour)) {
			continue
		}

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
				return err
			}
			continue
		}

		// create destination path
		vodsPath := filepath.Join("/var/www/media", "vods")
		segmentsPath := filepath.Join(vodsPath, newVod.Filename+"-segments")
		if err := os.MkdirAll(segmentsPath, 0755); err != nil && !os.IsExist(err) {
			return err
		}

		// get m3u8 playlist from twitch
		m3u8Url, err := external_apis.BuildDownloadURL(vod.ID, true)
		if err != nil {
			os.RemoveAll(segmentsPath)
			return err
		}

		// pass the m3u8 to ffmpeg to create .ts segments
		if err := createSegmentsfromURL(m3u8Url, segmentsPath, newVod.Filename, "vod"); err != nil {
			os.RemoveAll(segmentsPath)
			return err
		}

		// get metadata from m3u8
		var m meta
		m.Filename = newVod.Filename
		if err := getMetadata(segmentsPath, &m); err != nil {
			os.RemoveAll(segmentsPath)
			return err
		}
		newVod.Duration = m.Duration
		newVod.Resolution = m.Resolution
		newVod.Fps = m.Fps
		newVod.Size = m.Size

		// create thumbnails ...
		if err := createThumbnails(vodsPath, newVod.Filename, newVod.Duration); err != nil {
			return err
		}

		// create vod in database
		if err := queries.AddNewVod(&newVod); err != nil {
			return err
		}

		vods_downloaded += 1
	}

	if vods_downloaded == 0 {
		var settings models.Settings
		settings.DateVodsUpdate = time.Now()
		if err := queries.PartiallyUpdateSettings(&settings); err != nil {
			return err
		}
	}

	logger.Debug.Printf("[cronjob] vods downloaded: %d", vods_downloaded)
	return nil
}

func DownloadClips() error {
	logger.Debug.Println("[cronjob] clip download started")
	var clips []external_apis.TwitchHelixClip
	if err := external_apis.TwitchGetHelixClips(&clips); err != nil {
		return err
	}

	clips_downloaded := 0
	clips_updated := 0
	for _, clip := range clips {
		// skip clip if created less then 24h ago (terms of service)
		if !clip.CreatedAt.Before(time.Now().Add(time.Duration(-24) * time.Hour)) {
			continue
		}

		if clip.ViewCount < 3 {
			continue
		}

		// update game
		var game models.Game
		if err := queries.GetOneGame(&game, clip.GameID); err != nil {
			game.UUID = clip.GameID
			if e := queries.AddNewGame(&game); e != nil {
				return e
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
				return e
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
				return e
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
			return err
		}

		// get clip url from twitch
		downloadURL, err := external_apis.BuildDownloadURL(clip.ID, false)
		if err != nil {
			os.RemoveAll(segmentsPath)
			return err
		}

		// pass the clip url to ffmpeg to create .ts segments
		if err := createSegmentsfromURL(downloadURL, segmentsPath, clip.ID, "clip"); err != nil {
			os.RemoveAll(segmentsPath)
			return err
		}

		// get metadata from clip url
		var m meta
		m.Filename = clip.ID
		if err := getMetadata(segmentsPath, &m); err != nil {
			os.RemoveAll(segmentsPath)
			return err
		}
		newClip.Duration = m.Duration
		newClip.Resolution = m.Resolution
		newClip.Size = m.Size

		// create thumbnails ...
		if err := createThumbnails(clipsPath, clip.ID, newClip.Duration); err != nil {
			return err
		}

		// create clip in database
		if err := queries.AddNewClip(&newClip); err != nil {
			return err
		}

		clips_downloaded += 1
	}

	logger.Debug.Printf("[cronjob] clips downloaded: %d", clips_downloaded)
	logger.Debug.Printf("[cronjob] clips updated: %d", clips_updated)
	return nil
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
