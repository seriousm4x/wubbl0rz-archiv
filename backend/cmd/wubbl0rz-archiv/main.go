package main

import (
	"os"
	"runtime"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/assets"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/hooks"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
	_ "github.com/seriousm4x/wubbl0rz-archiv/internal/migrations"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/routes"
	"github.com/spf13/cobra"
)

func main() {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// probably ran with go run
		if runtime.GOOS == "windows" {
			assets.ArchiveDir = "Z:\\Archiv\\media"
		} else {
			assets.ArchiveDir = "/mnt/nas/Archiv/media"
		}
	} else {
		// probably ran with go build
		if runtime.GOOS == "windows" {
			assets.ArchiveDir = "Z:\\Archiv\\media"
		} else {
			assets.ArchiveDir = "/var/www/media/"
		}
	}

	app := pocketbase.New()

	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "createPreview",
		Short: "Create video preview for ids",
		Long:  "Create the hover webm/mp4 for given ids",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			assets.CreatePreview(app, args)
		},
	})

	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "createThumbnail",
		Short: "Create video thumbnail for ids",
		Long:  "Create webp thumbnails for given ids",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			assets.CreateThumbnail(app, args)
		},
	})

	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "createSprites",
		Short: "Create video sprites for ids",
		Long:  "Create video sprites for given ids",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			assets.CreateSprites(app, args)
		},
	})

	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "createPreviewThumbnailSprites",
		Short: "Create all assets for given ids",
		Long:  "Same as running createPreview, createThumbnail and createSprites but in one command",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			assets.CreatePreviewThumbnailsSprites(app, args)
		},
	})

	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "createAllAssets",
		Short: "Create all assets for the archive",
		Long:  "Creates assets for the entire archive forcefully. (Takes long time)",
		Run: func(cmd *cobra.Command, args []string) {
			assets.CreateAllAssets(app)
		},
	})

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Dir:         "internal/migrations",
		Automigrate: true,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(assets.ArchiveDir), false))
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/download/:type/:id", func(c echo.Context) error {
			return routes.Download(app, c)
		})
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/stats", func(c echo.Context) error {
			return routes.Stats(app, c)
		})
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		go func() {
			if err := hooks.InitBackend(app); err != nil {
				logger.Fatal.Fatalln(err)
			}
		}()
		return nil
	})

	// regenerate thumbnail if custom thumb has changed
	app.OnModelAfterUpdate("vod", "clip").Add(func(e *core.ModelEvent) error {
		oldRecord := e.Model.(*models.Record).OriginalCopy()
		oldCustomThumb := oldRecord.GetString("custom_thumbnail")
		newCustomThumb := e.Model.(*models.Record).GetString("custom_thumbnail")
		if oldCustomThumb != newCustomThumb {
			if err := assets.CreateThumbnail(app, []string{e.Model.GetId()}); err != nil {
				return err
			}
		}
		return nil
	})

	if err := app.Start(); err != nil {
		logger.Fatal.Fatalln(err)
	}
}
