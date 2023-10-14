package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/filesystem"
	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/logger"
)

func recreate(p string, id string) error {
	var m filesystem.Meta
	m.Filename = id

	if err := filesystem.GetMetadata(filepath.Join(p, id+"-segments"), &m); err != nil {
		return err
	}

	if err := filesystem.CreateThumbnails(p, id, m.Duration); err != nil {
		return err
	}

	return nil
}

func getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		return 0, 0, err
	}

	image, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, err
	}
	bounds := image.Bounds()
	return bounds.Max.X, bounds.Max.Y, nil
}

func main() {
	var files []string

	pathPtr := flag.String("path", "", "Path to the vods/clips base dir")
	flag.Parse()

	if _, err := os.Stat(*pathPtr); errors.Is(err, os.ErrNotExist) {
		logger.Error.Panicln(*pathPtr, "doesn't exist")
	}

	// find ids
	err := filepath.Walk(*pathPtr, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error.Panicln(err)
		}

		if info.IsDir() && strings.HasSuffix(path, "-segments") {
			filename := strings.Split(filepath.Base(path), "-segments")[0]
			files = append(files, filename)
		}

		return nil
	})

	if err != nil {
		logger.Error.Panicln(err)
	}

	type Thumbnail struct {
		Filename string
		Width    int
		Height   int
	}
	thumbnails := []Thumbnail{}

	for i, id := range files {
		logger.Info.Println(fmt.Sprintf("%d of %d: %s", i+1, len(files), id))

		thumbnails = []Thumbnail{}
		thumbnails = append(thumbnails, Thumbnail{Filename: "-sm.jpg", Width: 256, Height: 144})
		thumbnails = append(thumbnails, Thumbnail{Filename: "-md.jpg", Width: 512, Height: 288})
		thumbnails = append(thumbnails, Thumbnail{Filename: "-lg.jpg", Width: 1600, Height: 900})
		thumbnails = append(thumbnails, Thumbnail{Filename: "-sm.avif", Width: 256, Height: 144})
		thumbnails = append(thumbnails, Thumbnail{Filename: "-md.avif", Width: 512, Height: 288})

		for _, thumb := range thumbnails {
			imgPath := filepath.Join(*pathPtr, id+thumb.Filename)

			stat, err := os.Stat(imgPath)
			if errors.Is(err, os.ErrNotExist) {
				logger.Info.Println("Recrete", id, "(doesn't exist)")
				if err := recreate(*pathPtr, id); err != nil {
					logger.Error.Panicln(err)
				}
				continue
			}

			if stat.IsDir() {
				continue
			}

			// go cant decode avif to get dimensions
			if strings.HasSuffix(imgPath, ".avif") {
				continue
			}

			width, height, err := getImageDimension(imgPath)
			if stat.Size() <= 8 || width != thumb.Width || height != thumb.Height || err != nil {
				logger.Info.Println("Recrete", id)
				if err := recreate(*pathPtr, id); err != nil {
					logger.Error.Panicln(err)
				}
			}
		}
	}
}
