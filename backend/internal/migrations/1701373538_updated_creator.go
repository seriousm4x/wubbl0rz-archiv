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

		collection, err := dao.FindCollectionByNameOrId("u6cxge211isaeir")
		if err != nil {
			return err
		}

		// add
		new_ttv_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "9lsegz2y",
			"name": "ttv_id",
			"type": "number",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"noDecimal": false
			}
		}`), new_ttv_id)
		collection.Schema.AddField(new_ttv_id)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("u6cxge211isaeir")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("9lsegz2y")

		return dao.SaveCollection(collection)
	})
}
