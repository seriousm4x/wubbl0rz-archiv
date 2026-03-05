package cronjobs

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

func RunVaccum(app core.App) error {
	logger.Debug.Println("[cronjob] running vacuum")

	if err := app.Vacuum(); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
