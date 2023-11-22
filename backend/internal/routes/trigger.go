package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/cronjobs"
)

func TriggerTwitchDownloads(c echo.Context) error {
	go cronjobs.RunTwitchDownloads(App)
	return c.NoContent(200)
}
