package hooks

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

// Import all environment variables to the database.
func ImportEnv(app *pocketbase.PocketBase) error {
	logger.Debug.Println("[hooks] importing env")

	ttvClientId := os.Getenv("TTV_CLIENT_ID")
	if ttvClientId == "" {
		err := fmt.Errorf("os env \"TTV_CLIENT_ID\" must not be empty")
		logger.Error.Println(err)
		return err
	}

	ttvClientSecret := os.Getenv("TTV_CLIENT_SECRET")
	if ttvClientSecret == "" {
		err := fmt.Errorf("os env \"TTV_CLIENT_SECRET\" must not be empty")
		logger.Error.Println(err)
		return err
	}

	var settings *core.Record

	settings, err := app.FindFirstRecordByFilter("settings", "id != ''")
	if err == sql.ErrNoRows {
		collection, err := app.FindCollectionByNameOrId("settings")
		if err != nil {
			logger.Error.Println(err)
			return err
		}
		settings = core.NewRecord(collection)
	} else if err != nil {
		logger.Error.Println(err)
		return err
	}

	settings.Set("ttv_client_id", ttvClientId)
	settings.Set("ttv_client_secret", ttvClientSecret)

	if err := app.Save(settings); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
