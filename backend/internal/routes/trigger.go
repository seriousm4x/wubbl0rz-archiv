package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/cronjobs"
)

func TriggerTwitchDownloads(app *pocketbase.PocketBase, c echo.Context) error {
	go cronjobs.RunTwitchDownloads(app)
	return c.NoContent(200)
}
