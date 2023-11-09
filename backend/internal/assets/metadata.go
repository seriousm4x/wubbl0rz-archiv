package assets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

type Meta struct {
	Filename   string
	Duration   int
	Resolution string
	Fps        float32
	Size       int
}

type FFProbe struct {
	Programs []struct {
		Streams []struct {
			Width      int
			Height     int
			RFrameRate string `json:"r_frame_rate"`
		}
	}
	Format struct {
		Duration string
	}
}

// Sums up all segment sizes
func getSegmentSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error.Println(err)
			return err
		}
		if !info.IsDir() && strings.Contains(info.Name(), ".ts") {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// Extract metadata with ffprobe
func GetMetadata(destPath string, m *Meta) error {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
		"program_stream=width,height,r_frame_rate:format=duration", "-of", "json",
		filepath.Join(destPath, m.Filename+".m3u8"))

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.Error.Println(err)
		return err
	}

	if stderr.String() != "" {
		logger.Error.Println(stderr)
		return errors.New(stderr.String())
	}

	var ffprobe FFProbe
	if err := json.Unmarshal(stdout.Bytes(), &ffprobe); err != nil {
		logger.Error.Println(err)
		return err
	}

	if len(ffprobe.Programs) == 0 {
		err := errors.New("ffprobe: no entries in json key 'programs'")
		logger.Error.Println(err)
		logger.Error.Println(stdout.String())
		return err
	} else if len(ffprobe.Programs[0].Streams) == 0 {
		err := errors.New("ffprobe: no entries in json key 'streams'")
		logger.Error.Println(err)
		logger.Error.Println(stdout.String())
		return err
	}

	// width, height
	width := ffprobe.Programs[0].Streams[0].Width
	height := ffprobe.Programs[0].Streams[0].Height

	// fps
	fpsFraction := strings.Split(ffprobe.Programs[0].Streams[0].RFrameRate, "/")
	fpsNumerator, err := strconv.ParseFloat(fpsFraction[0], 64)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	fpsDenominator, err := strconv.ParseFloat(fpsFraction[1], 64)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	fps := fpsNumerator / fpsDenominator

	// duration
	duration, err := strconv.ParseFloat(ffprobe.Format.Duration, 64)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	m.Duration = int(math.Round(duration))
	m.Resolution = fmt.Sprintf("%dx%d", width, height)
	m.Fps = float32(fps)

	// get filesize
	size, err := getSegmentSize(destPath)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	m.Size = int(size)

	return nil
}
