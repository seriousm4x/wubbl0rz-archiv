package routes

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/cronjobs"
)

func TriggerTwitchDownloads(e *core.RequestEvent) error {
	go cronjobs.RunTwitchDownloads(App)
	return e.NoContent(200)
}
