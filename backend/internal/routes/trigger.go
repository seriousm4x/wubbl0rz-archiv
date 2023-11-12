package routes

import (
	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/cronjobs"
)

func TriggerTwitchDownloads(app *pocketbase.PocketBase) error {
	go cronjobs.RunTwitchDownloads(app)
	return nil
}
