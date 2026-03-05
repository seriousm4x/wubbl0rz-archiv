package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("cyuvby822bxphoz")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("dlxgq67c")

		// remove field
		collection.Fields.RemoveById("gx1q3pdo")

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("cyuvby822bxphoz")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
			"hidden": false,
			"id": "dlxgq67c",
			"maxSize": 2000000,
			"name": "yt_client_secret",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "json"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
			"hidden": false,
			"id": "gx1q3pdo",
			"maxSize": 2000000,
			"name": "yt_bearer_token",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "json"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
