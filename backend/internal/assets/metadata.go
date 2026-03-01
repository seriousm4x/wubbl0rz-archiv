package assets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
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
	Streams []struct {
		Width        int
		Height       int
		AvgFrameRate string `json:"avg_frame_rate"`
	}
	Format struct {
		Duration string
	}
}

// Extract metadata with ffprobe
func GetMetadata(mp4 string, m *Meta) error {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
		"stream=width,height,avg_frame_rate:format=duration", "-of", "json",
		mp4)

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

	if len(ffprobe.Streams) == 0 {
		err := errors.New("ffprobe: no entries in json key 'streams'")
		logger.Error.Println(err)
		logger.Error.Println(stdout.String())
		return err
	}

	// width, height
	width := ffprobe.Streams[0].Width
	height := ffprobe.Streams[0].Height

	// fps
	fpsFraction := strings.Split(ffprobe.Streams[0].AvgFrameRate, "/")
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
	info, err := os.Stat(mp4)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	m.Size = int(info.Size())

	return nil
}
