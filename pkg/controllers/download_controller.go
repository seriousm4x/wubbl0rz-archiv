package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"
	"github.com/gin-gonic/gin"
)

func SendStream(c *gin.Context) {
	media_type := c.Param("type")
	media_path := "/var/www/media/"
	uuid := c.Param("uuid")

	var vod models.Vod
	var clip models.Clip
	var filename string
	var title string

	if media_type == "vods" {
		if err := queries.GetOneVod(&vod, uuid, true); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": true,
				"msg":   "Vod not found",
			})
			return
		}
		filename = vod.Filename
		title = vod.Title
	} else if media_type == "clips" {
		if err := queries.GetOneClip(&clip, uuid); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": true,
				"msg":   "Clip not found",
			})
			return
		}
		filename = clip.Filename
		title = clip.Title
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   fmt.Sprintf("%s is not a valid media type", media_type),
		})
		return
	}

	input_file := filepath.Join(media_path, media_type, filename+"-segments", filename+".m3u8")

	cmd := exec.Command("ffmpeg",
		"-i", input_file,
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		"-movflags", "frag_keyframe+empty_moov",
		"-f", "mp4",
		"-")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	data := make([]byte, 1024)
	for {
		c.Writer.Header().Set("Content-Type", "video/mp4")
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s.mp4", filepath.Clean(title)))
		n, err := stdout.Read(data)
		if err == io.EOF {
			// end of file
			break
		}
		if n > 0 {
			valid_data := data[:n]
			if _, err := c.Writer.Write(valid_data); err != nil {
				// client disconnected. kill ffmpeg
				cmd.Process.Kill()
				break
			}
		}
	}
}
