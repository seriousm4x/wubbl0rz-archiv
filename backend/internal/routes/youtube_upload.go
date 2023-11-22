package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/seriousm4x/wubbl0rz-archiv/external"
)

func YoutubeUpload(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "id must not be empty",
		})
	}

	go external.YoutubeUpload(App, id)

	return c.NoContent(http.StatusNoContent)
}
