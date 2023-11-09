package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "9wu419qp30znepf",
			"created": "2023-11-09 19:56:30.938Z",
			"updated": "2023-11-09 19:56:30.938Z",
			"name": "vod",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "vbdiu1fg",
					"name": "title",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "yi7rcxpk",
					"name": "duration",
					"type": "number",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"noDecimal": false
					}
				},
				{
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
				},
				{
					"system": false,
					"id": "rbd4cxr7",
					"name": "viewcount",
					"type": "number",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"noDecimal": false
					}
				},
				{
					"system": false,
					"id": "5cdrv0os",
					"name": "filename",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "5ay4apy8",
					"name": "resolution",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "fu9c5gef",
					"name": "fps",
					"type": "number",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"noDecimal": false
					}
				},
				{
					"system": false,
					"id": "smojzwjt",
					"name": "size",
					"type": "number",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"noDecimal": false
					}
				},
				{
					"system": false,
					"id": "rv2bfwwv",
					"name": "publish",
					"type": "bool",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "wprkbol0",
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
				}
			],
			"indexes": [],
			"listRule": "publish = true",
			"viewRule": "publish = true",
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("9wu419qp30znepf")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
