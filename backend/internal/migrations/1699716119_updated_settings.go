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

		collection, err := dao.FindCollectionByNameOrId("cyuvby822bxphoz")
		if err != nil {
			return err
		}

		// add
		new_yt_client_secret := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "dlxgq67c",
			"name": "yt_client_secret",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_yt_client_secret)
		collection.Schema.AddField(new_yt_client_secret)

		// add
		new_yt_bearer_token := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "gx1q3pdo",
			"name": "yt_bearer_token",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_yt_bearer_token)
		collection.Schema.AddField(new_yt_bearer_token)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("cyuvby822bxphoz")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("dlxgq67c")

		// remove
		collection.Schema.RemoveField("gx1q3pdo")

		return dao.SaveCollection(collection)
	})
}
