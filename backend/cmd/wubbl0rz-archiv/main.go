package main

import (
	"os"
	"runtime"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
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
		switch runtime.GOOS {
		case "windows":
			assets.ArchiveDir = "Z:\\Archiv\\media"
		case "darwin":
			assets.ArchiveDir = "/Volumes/nas/Archiv/media"
		default:
			assets.ArchiveDir = "/mnt/nas/Archiv/media"
		}
	} else {
		// probably ran with go build
		switch runtime.GOOS {
		case "windows":
			assets.ArchiveDir = "Z:\\Archiv\\media"
		case "darwin":
			assets.ArchiveDir = "/Volumes/nas/Archiv/media"
		default:
			assets.ArchiveDir = "/var/www/media"
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

	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "createAudioOnly",
		Short: "Create an audio only ogg file",
		Long:  "Creates an audio only ogg file for vods",
		Run: func(cmd *cobra.Command, args []string) {
			assets.CreateAudioOnly(app, args)
		},
	})

	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "refreshMetadata",
		Short: "Refresh metadata for vods and clips",
		Long:  "Refresh duration, resolution, fps and size in database for vods and clips.",
		Run: func(cmd *cobra.Command, args []string) {
			assets.RefreshMetadata(app)
		},
	})

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Dir:         "internal/migrations",
		Automigrate: true,
	})

	// public routes
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS(assets.ArchiveDir), false))

		// route for statistics
		se.Router.GET("/stats", func(e *core.RequestEvent) error {
			return routes.Stats(app, e)
		})

		return se.Next()
	})

	// auth routes
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// route for triggering twitch downloads
		se.Router.GET("/wubbl0rz/trigger/downloads",
			routes.TriggerTwitchDownloads).Bind(apis.RequireAuth("_superusers"))
		return se.Next()
	})

	// init backend once on start
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		go func() {
			if err := hooks.InitBackend(app); err != nil {
				logger.Fatal.Fatalln(err)
			}
		}()
		return se.Next()
	})

	// regenerate thumbnail if custom thumb has changed
	app.OnRecordAfterUpdateSuccess("vod", "clip").BindFunc(func(e *core.RecordEvent) error {
		oldRecord := e.Record.Original()
		oldCustomThumb := oldRecord.GetString("custom_thumbnail")
		newCustomThumb := e.Record.GetString("custom_thumbnail")
		if oldCustomThumb != newCustomThumb {
			if err := assets.CreateThumbnail(app, []string{e.Record.Id}); err != nil {
				return err
			}
		}
		return nil
	})

	if err := app.Start(); err != nil {
		logger.Fatal.Fatalln(err)
	}
}
