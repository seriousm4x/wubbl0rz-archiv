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

		collection, err := dao.FindCollectionByNameOrId("9wu419qp30znepf")
		if err != nil {
			return err
		}

		// update
		edit_youtube_upload := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "d5keuz5s",
			"name": "youtube_upload",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"pending",
					"done"
				]
			}
		}`), edit_youtube_upload)
		collection.Schema.AddField(edit_youtube_upload)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("9wu419qp30znepf")
		if err != nil {
			return err
		}

		// update
		edit_youtube_upload := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "d5keuz5s",
			"name": "youtube_upload",
			"type": "select",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"none",
					"pending",
					"done"
				]
			}
		}`), edit_youtube_upload)
		collection.Schema.AddField(edit_youtube_upload)

		return dao.SaveCollection(collection)
	})
}
