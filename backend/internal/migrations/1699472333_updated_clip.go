package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ze75ahzrvxoq372")
		if err != nil {
			return err
		}

		// update
		edit_custom_thumbnail := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ur9iha8g",
			"name": "custom_thumbnail",
			"type": "file",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"maxSize": 20970000,
				"mimeTypes": [
					"image/jpeg",
					"image/png",
					"image/svg+xml",
					"image/gif",
					"image/webp",
					"image/avif",
					"image/tiff"
				],
				"thumbs": [],
				"protected": false
			}
		}`), edit_custom_thumbnail)
		collection.Schema.AddField(edit_custom_thumbnail)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ze75ahzrvxoq372")
		if err != nil {
			return err
		}

		// update
		edit_custom_thumbnail := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ur9iha8g",
			"name": "custom_thumbnail",
			"type": "file",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"maxSize": 5242880,
				"mimeTypes": [
					"image/jpeg",
					"image/png",
					"image/svg+xml",
					"image/gif",
					"image/webp",
					"image/avif",
					"image/tiff"
				],
				"thumbs": [],
				"protected": false
			}
		}`), edit_custom_thumbnail)
		collection.Schema.AddField(edit_custom_thumbnail)

		return dao.SaveCollection(collection)
	})
}
