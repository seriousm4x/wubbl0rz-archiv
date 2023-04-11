package queries

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/database"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"gorm.io/gorm"
)

func GetAllVods(v *[]models.Vod, query models.Vod, pagination Pagination, o string) (*Pagination, error) {
	if o == "" {
		o = "date desc"
	}
	result := database.DB.Model(&query).Omit("transcript")
	if query.Title != "" {
		// if title is given, do case insensitive search in title string
		result = result.Where("position(LOWER(?) in LOWER(title))>0", query.Title)
	} else {
		// else search exact query match
		result = result.Where(query)
	}
	result = result.Where("publish = ?", true).Order(o).Count(&pagination.TotalRows).Scopes(Paginate(&pagination, database.DB)).Find(v)
	if result.RowsAffected == 0 {
		return &pagination, errors.New("not found")
	}
	return &pagination, nil
}

func AddNewVod(v *models.Vod) error {
	if v.UUID == "" {
		b := make([]byte, 3)
		for {
			rand.Read(b)
			uuid := hex.EncodeToString(b)
			if err := GetOneVod(&models.Vod{}, uuid, false); err != nil {
				v.UUID = uuid
				break
			}
		}
	}

	if err := database.DB.Create(v).Error; err != nil {
		return err
	}

	// update vod timestamp
	var settings models.Settings
	settings.DateVodsUpdate = time.Now()
	if err := PartiallyUpdateSettings(&settings); err != nil {
		return err
	}

	return nil
}

func GetOneVod(v *models.Vod, uuid string, onlyPublic bool) error {
	var result *gorm.DB
	if onlyPublic {
		result = database.DB.Where("uuid = ?", uuid).Where("publish = ?", true).Preload("Clips.Creator").Preload("Clips.Game")
	} else {
		result = database.DB.Where("uuid = ?", uuid)
	}
	result = result.Find(v)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetVodByFilename(v *models.Vod, filename string) error {
	result := database.DB.Omit("transcript").Where("filename = ?", filename).Find(v)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetVodsByUUID(v *[]models.Vod, uuids []string) error {
	result := database.DB.Omit("transcript").Where("uuid IN ?", uuids).Where("publish = ?", true).Find(v)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetVodsByYear(v *[]models.Vod, year string) error {
	result := database.DB.Model(&v).Omit("transcript").Where("date_part('year', date) = ?", year).Where("publish = ?", true).Order("date desc").Find(v)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func PatchVod(changes map[string]interface{}, uuid string) error {
	var vod models.Vod
	if err := GetOneVod(&vod, uuid, false); err != nil {
		return errors.New("vod not found")
	}

	if err := database.DB.Model(&vod).Where("uuid = ?", uuid).Updates(changes).Error; err != nil {
		return errors.New("update failed")
	}

	// update vod timestamp
	var settings models.Settings
	settings.DateVodsUpdate = time.Now()
	if err := PartiallyUpdateSettings(&settings); err != nil {
		return err
	}

	return nil
}

func DeleteVod(v *models.Vod, uuid string) error {
	database.DB.Where("uuid = ?", uuid).Delete(v)
	return nil
}
