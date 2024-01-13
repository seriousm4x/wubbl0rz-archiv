package assets

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

var ArchiveDir string

// Takes seconds and formats it as VTT timecode
func secondsToWebVTTTimecode(seconds int) string {
	hours := seconds / 3600
	remainingSeconds := seconds % 3600
	minutes := remainingSeconds / 60
	seconds = remainingSeconds % 60

	formattedHours := fmt.Sprintf("%02d", hours)
	formattedMinutes := fmt.Sprintf("%02d", minutes)
	formattedSeconds := fmt.Sprintf("%02d", seconds)
	formattedMilliseconds := "000"

	webVTTTimecode := fmt.Sprintf("%s:%s:%s.%s", formattedHours, formattedMinutes, formattedSeconds, formattedMilliseconds)

	return webVTTTimecode
}

// Taskes an id and returns the vod/clip record
func findVodOrClip(app *pocketbase.PocketBase, id string) (*models.Record, error) {
	var record *models.Record
	var err error
	record, err = app.Dao().FindRecordById("vod", id)
	if err != nil {
		record, err = app.Dao().FindRecordById("clip", id)
		if err != nil {
			logger.Error.Printf("Record for id %s not found: %v", id, err)
			return record, err
		}
	}
	return record, err
}

// Create the video preview files for multiple ids (files for video hover)
func CreatePreview(app *pocketbase.PocketBase, ids []string) error {
	for i, id := range ids {
		logger.Debug.Printf("[%d of %d] Recreating video preview for id \"%s\" ", i+1, len(ids), id)
		record, err := findVodOrClip(app, id)
		if err != nil {
			return err
		}

		mediaType := fmt.Sprintf("%ss", record.Collection().Name)
		m3u8 := filepath.Join(ArchiveDir, mediaType, fmt.Sprintf("%s%s", record.GetString("filename"), "-segments"), record.GetString("filename")+".m3u8")
		duration := record.GetInt("duration")

		var cmd *exec.Cmd
		outputWebm := filepath.Join(ArchiveDir, mediaType, record.GetString("filename")+"-preview.webm")
		outputMp4 := filepath.Join(ArchiveDir, mediaType, record.GetString("filename")+"-preview.mp4")

		if mediaType == "vods" && duration > 60 {
			// create video preview with 3 segments
			seekPoint1 := duration / 4
			seekPoint2 := 2 * seekPoint1
			seekPoint3 := 3 * seekPoint1
			segmentLength := "3"

			// this takes the vod/clip as an input and defines 3 seek ranges (seekPointX + segmentLength) which will be cut together to create
			// a better summarization of the vod
			// outputs 2 videos: one vp9 webm and one h265 mp4
			cmd = exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", strconv.Itoa(seekPoint1), "-t", segmentLength, "-i", m3u8, "-ss", strconv.Itoa(seekPoint2), "-t", segmentLength, "-i", m3u8, "-ss", strconv.Itoa(seekPoint3), "-t", segmentLength, "-i", m3u8, "-filter_complex", "[0:v]scale=-2:720[0v];[1:v]scale=-2:720[1v];[2:v]scale=-2:720[2v];[0v][1v][2v]concat=n=3:v=1,split=2[out1][out2]", "-map", "[out1]", "-c:v:0", "libvpx-vp9", "-crf", "38", "-b:v", "0", "-r", "25", "-an", "-y", outputWebm, "-map", "[out2]", "-c:v:1", "libx265", "-crf", "28", "-r", "25", "-an", "-y", outputMp4)
		} else {
			// just take first 4 seconds as preview
			// outputs 2 videos: one vp9 webm and one h265 mp4
			cmd = exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-t", "4", "-i", m3u8, "-filter_complex", "[0:v]scale=-2:720[0v];[0v]split=2[out1][out2]", "-map", "[out1]", "-c:v:0", "libvpx-vp9", "-crf", "38", "-b:v", "0", "-r", "25", "-an", "-y", outputWebm, "-map", "[out2]", "-c:v:1", "libx265", "-crf", "28", "-r", "25", "-an", "-y", outputMp4)
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			logger.Error.Printf("%v", stderr.String())
			logger.Error.Println(err.Error())
			return err
		}
	}

	return nil
}

// Create the video thumbnails for multiple ids
func CreateThumbnail(app *pocketbase.PocketBase, ids []string) error {
	for i, id := range ids {
		logger.Debug.Printf("[%d of %d] Recreating video thumbnail for id \"%s\" ", i+1, len(ids), id)
		record, err := findVodOrClip(app, id)
		if err != nil {
			return err
		}

		var cmd *exec.Cmd
		customThumbnail := record.GetString("custom_thumbnail")
		mediaType := fmt.Sprintf("%ss", record.Collection().Name)

		// create sm, md and lg webp thumbs thumb
		outputSm := filepath.Join(ArchiveDir, mediaType, record.GetString("filename")+"-sm.webp")
		outputMd := filepath.Join(ArchiveDir, mediaType, record.GetString("filename")+"-md.webp")
		outputLg := filepath.Join(ArchiveDir, mediaType, record.GetString("filename")+"-lg.webp")
		compression := "0"
		quality := "85"

		if customThumbnail == "" {
			// create thumb from m3u8
			m3u8 := filepath.Join(ArchiveDir, mediaType, fmt.Sprintf("%s%s", record.GetString("filename"), "-segments"), record.GetString("filename")+".m3u8")
			duration := record.GetInt("duration")
			var timecode_framegrab string
			if duration <= 10 {
				timecode_framegrab = "0"
			} else {
				timecode_framegrab = fmt.Sprintf("%d", int(duration/2))
			}
			cmd = exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
				"-filter_complex", "[0:v]split=3[frame1][frame2][frame3];[frame1]scale=512:-2[outputSm];[frame2]scale=768:-2[outputMd];[frame3]scale=1536:-2[outputLg]", "-map", "[outputSm]", "-c:v", "libwebp", "-frames", "1", "-compression_level", compression, "-quality", quality, "-y", outputSm, "-map", "[outputMd]", "-c:v", "libwebp", "-frames", "1", "-compression_level", compression, "-quality", quality, "-y", outputMd, "-map", "[outputLg]", "-c:v", "libwebp", "-frames", "1", "-compression_level", compression, "-quality", quality, "-y", outputLg)
		} else {
			// create thumb from custom image
			input := path.Join(app.DataDir(), "storage", record.Collection().Id, record.Id, customThumbnail)
			cmd = exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-i", input,
				"-filter_complex", "[0:v]split=3[frame1][frame2][frame3];[frame1]scale=512:-2[outputSm];[frame2]scale=768:-2[outputMd];[frame3]scale=1536:-2[outputLg]", "-map", "[outputSm]", "-c:v", "libwebp", "-frames", "1", "-compression_level", compression, "-quality", quality, "-y", outputSm, "-map", "[outputMd]", "-c:v", "libwebp", "-frames", "1", "-compression_level", compression, "-quality", quality, "-y", outputMd, "-map", "[outputLg]", "-c:v", "libwebp", "-frames", "1", "-compression_level", compression, "-quality", quality, "-y", outputLg)
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			logger.Error.Printf("%v", stderr.String())
			logger.Error.Println(err)
			return err
		}

	}

	return nil
}

// Create the video sprites for multiple ids (hover over player slider)
func CreateSprites(app *pocketbase.PocketBase, ids []string) error {
	for i, id := range ids {
		logger.Debug.Printf("[%d of %d] Recreating video sprites for id \"%s\" ", i+1, len(ids), id)
		record, err := findVodOrClip(app, id)
		if err != nil {
			return err
		}

		mediaType := fmt.Sprintf("%ss", record.Collection().Name)
		m3u8 := filepath.Join(ArchiveDir, mediaType, fmt.Sprintf("%s%s", record.GetString("filename"), "-segments"), record.GetString("filename")+".m3u8")

		outDir := filepath.Join(ArchiveDir, mediaType, record.GetString("filename")+"-sprites")
		if err := os.MkdirAll(outDir, 700); err != nil {
			logger.Error.Println(err)
			return err
		}

		outputSprites := filepath.Join(outDir, record.GetString("filename")+"_%03d.webp")
		cmd := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-skip_frame", "nokey", "-i", m3u8, "-vf", "fps=1/30,scale=192:-2,tile=10x10", "-c:v", "libwebp", "-compression_level", "0", "-quality", "85", "-y", outputSprites)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			logger.Error.Printf("%v", stderr.String())
			logger.Error.Println(err)
			return err
		}

		// create vtt
		duration := record.GetFloat("duration") // 12265.0
		secondsPerFrame := 30.0
		imagesCount := math.Ceil(duration / secondsPerFrame)

		webvtt := "WEBVTT\n"
		imgWidth := 192
		imgHeight := 108
		currentX, currentRow := 0, 0
		lastFilenameCount := 0

		for index := range make([]int, int(imagesCount)) {
			filenameCount := int(math.Floor(float64(index)/100)) + 1

			// reset positions on new file
			if lastFilenameCount != filenameCount {
				currentX, currentRow = 0, 0
			}

			// end of line reached, reset to new line
			if currentX >= 10 {
				currentX = 0
				currentRow++
			}
			xPos := (currentX * imgWidth)
			yPos := (currentRow * imgHeight)

			// count up x
			if currentX >= 0 && currentX <= 10 {
				xPos = (currentX * imgWidth)
			}

			// count up y
			if currentRow >= 0 && currentRow < 10 {
				yPos = (currentRow * imgHeight)
			}

			secondsStart := index * int(secondsPerFrame)
			secondsEnd := index*int(secondsPerFrame) + int(secondsPerFrame)
			timecodeStart := secondsToWebVTTTimecode(secondsStart)
			timecodeEnd := secondsToWebVTTTimecode(secondsEnd)

			str := fmt.Sprintf("\n%s --> %s\n%s_%s.webp#xywh=%d,%d,%d,%d\n",
				timecodeStart, timecodeEnd, record.GetString("filename"), fmt.Sprintf("%03d", filenameCount),
				xPos, yPos, imgWidth, imgHeight)
			webvtt += str

			currentX = currentX + 1
			lastFilenameCount = filenameCount
		}

		outputVTT := filepath.Join(outDir, record.GetString("filename")+".vtt")

		if err := os.WriteFile(outputVTT, []byte(webvtt), 0644); err != nil {
			logger.Error.Println(err)
			return err
		}
	}
	return nil
}

// Creates video preview, thumbnails and sprites
func CreatePreviewThumbnailsSprites(app *pocketbase.PocketBase, ids []string) error {
	// create wait group to run jobs parallel
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := CreatePreview(app, ids); err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := CreateThumbnail(app, ids); err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := CreateSprites(app, ids); err != nil {
			return
		}
	}()

	wg.Wait()

	return nil
}

// Recreates assets for the entire archive forcefully. (Takes long time)
func CreateAllAssets(app *pocketbase.PocketBase) error {
	// vods
	allVods, err := app.Dao().FindRecordsByExpr("vod")
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	var vodIds []string
	for _, vod := range allVods {
		vodIds = append(vodIds, vod.Id)
	}

	if err := CreatePreviewThumbnailsSprites(app, vodIds); err != nil {
		return err
	}

	// clips
	allClips, err := app.Dao().FindRecordsByExpr("clip")
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	var clipIds []string
	for _, clip := range allClips {
		clipIds = append(clipIds, clip.Id)
	}

	if err := CreatePreviewThumbnailsSprites(app, clipIds); err != nil {
		return err
	}

	return nil
}
