package routes

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/external"
)

func YoutubeUpload(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	if id == "" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"error": "id must not be empty",
		})
	}

	go external.YoutubeUpload(App, id)

	return e.NoContent(http.StatusNoContent)
}
