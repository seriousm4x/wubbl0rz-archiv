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
		edit_date := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "u4zrms8a",
			"name": "date",
			"type": "date",
			"required": false,
			"presentable": true,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), edit_date)
		collection.Schema.AddField(edit_date)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("9wu419qp30znepf")
		if err != nil {
			return err
		}

		// update
		edit_date := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "u4zrms8a",
			"name": "date",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), edit_date)
		collection.Schema.AddField(edit_date)

		return dao.SaveCollection(collection)
	})
}
