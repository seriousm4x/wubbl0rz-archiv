package filesystem

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Meta struct {
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

func GetMetadata(destPath string, m *Meta) error {
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

	// width, height
	width := strings.TrimSpace(splittedStdout[0])
	height := strings.TrimSpace(splittedStdout[1])

	// fps
	fpsFraction := strings.Split(strings.TrimSpace(splittedStdout[2]), "/")
	fpsNumerator, err := strconv.ParseFloat(fpsFraction[0], 64)
	if err != nil {
		return err
	}
	fpsDenominator, err := strconv.ParseFloat(fpsFraction[1], 64)
	if err != nil {
		return err
	}
	fps := fpsNumerator / fpsDenominator

	// duration
	duration, err := strconv.ParseFloat(strings.TrimSpace(splittedStdout[3]), 64)
	if err != nil {
		return err
	}

	// check empty values
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

func CreateThumbnails(destPath string, filename string, duration int) error {
	m3u8 := filepath.Join(destPath, filename+"-segments", filename+".m3u8")

	var timecode_framegrab string
	if duration <= 10 {
		timecode_framegrab = "0"
	} else {
		timecode_framegrab = fmt.Sprintf("%d", int(duration/2))
	}

	// create lossless source png
	src_png := filepath.Join(destPath, filename+"-source.png")
	cmd := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab,
		"-i", m3u8, "-vframes", "1", "-f", "image2", "-y", src_png)
	if err := cmd.Run(); err != nil {
		return err
	}

	// define final jpg thumbnails
	type JpgThumbnail struct {
		Filename string
		Width    int
		Height   int
	}
	jpgThumbnails := []JpgThumbnail{}
	jpgThumbnails = append(jpgThumbnails, JpgThumbnail{Filename: "-sm.jpg", Width: 256, Height: 144})
	jpgThumbnails = append(jpgThumbnails, JpgThumbnail{Filename: "-md.jpg", Width: 512, Height: 288})
	jpgThumbnails = append(jpgThumbnails, JpgThumbnail{Filename: "-lg.jpg", Width: 1600, Height: 900})

	// define final avif thumbnails
	type AvifThumbnail struct {
		Filename string
		Width    int
		Height   int
	}
	avifThumbnails := []AvifThumbnail{}
	avifThumbnails = append(avifThumbnails, AvifThumbnail{Filename: "-sm.avif", Width: 256, Height: 144})
	avifThumbnails = append(avifThumbnails, AvifThumbnail{Filename: "-md.avif", Width: 512, Height: 288})

	// encode jpg and avif thumbnails with ffmpeg
	for _, thumb := range jpgThumbnails {
		cmd = exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-i", src_png,
			"-vf", fmt.Sprintf("scale=%d:%d", thumb.Width, thumb.Height), "-y", filepath.Join(destPath, filename+thumb.Filename))
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	for _, thumb := range avifThumbnails {
		cmd = exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-i", src_png, "-c:v", "libaom-av1",
			"-still-picture", "1", "-vf", fmt.Sprintf("scale=%d:%d", thumb.Width, thumb.Height), "-y",
			filepath.Join(destPath, filename+thumb.Filename))
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// remove source png
	if err := os.Remove(src_png); err != nil {
		return err
	}

	// animated webp
	animated_webp := filepath.Join(destPath, filename+"-preview.webp")
	cmd = exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-i", m3u8, "-ss",
		timecode_framegrab, "-c:v", "libwebp", "-vf", "scale=256:-1,fps=fps=15", "-lossless",
		"0", "-compression_level", "3", "-q:v", "70", "-loop", "0", "-preset", "picture",
		"-an", "-vsync", "0", "-t", "4", "-y", animated_webp)
	if err := cmd.Run(); err != nil {
		return err
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
