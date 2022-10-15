package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/AgileProggers/archiv-backend-go/pkg/filesystem"
	"github.com/AgileProggers/archiv-backend-go/pkg/logger"
)

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

	for _, id := range files {
		logger.Debug.Println(id)

		var m filesystem.Meta
		m.Filename = id

		if err := filesystem.GetMetadata(filepath.Join(*pathPtr, id+"-segments"), &m); err != nil {
			logger.Error.Panicln(err)
		}

		if err := filesystem.CreateThumbnails(*pathPtr, id, m.Duration); err != nil {
			logger.Error.Panicln(err)
		}

	}

	if err != nil {
		logger.Error.Panicln(err)
	}
}
