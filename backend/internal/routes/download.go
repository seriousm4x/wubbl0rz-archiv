package routes

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/assets"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// Route to download vods and clips
func Download(app *pocketbase.PocketBase, c echo.Context) error {
	media_type := c.PathParam("type")
	id := c.PathParam("id")

	var title string
	var filename string

	switch {
	case media_type == "vods":
		vod, err := app.Dao().FindRecordById("vod", id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]any{
				"message": fmt.Sprintf("No vod with id '%s' found.", id),
			})
		}
		title = vod.GetString("title")
		filename = vod.GetString("filename")
	case media_type == "clips":
		clip, err := app.Dao().FindRecordById("clip", id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]any{
				"message": fmt.Sprintf("No clip with id '%s' found.", id),
			})
		}
		title = clip.GetString("title")
		filename = clip.GetString("filename")
	default:
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": fmt.Sprintf("'%s' is not a valid media type. Only 'vods' and 'clips' are valid.", media_type),
		})
	}

	cmd := exec.Command("ffmpeg",
		"-i", filepath.Join(assets.ArchiveDir, media_type, filename+"-segments", filename+".m3u8"),
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		"-movflags", "frag_keyframe+empty_moov",
		"-f", "mp4",
		"-")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to prepare process",
		})
	}
	defer stdout.Close()

	if err = cmd.Start(); err != nil {
		logger.Error.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to execute process",
		})
	}

	buf := make([]byte, 1024*1024*10) // 10 MB
	for {
		c.Response().Header().Set(echo.HeaderContentType, "video/mp4")
		c.Response().Header().Add(echo.HeaderContentDisposition,
			fmt.Sprintf("attachment; filename=%s.mp4", filepath.Clean(title)))

		n, err := stdout.Read(buf)
		if err == io.EOF {
			break
		}

		if n > 0 {
			chunk := buf[:n]
			if _, err := c.Response().Writer.Write(chunk); err != nil {
				// client disconnected or other error. kill ffmpeg
				cmd.Process.Kill()
				break
			}
			c.Response().Flush()
		}
	}

	if err := cmd.Wait(); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
