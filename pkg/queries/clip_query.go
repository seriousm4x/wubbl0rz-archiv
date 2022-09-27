package queries

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/database"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
)

func GetAllClips(c *[]models.Clip, query models.Clip, pagination Pagination, o string, date_from time.Time, date_to time.Time) (*Pagination, error) {
	if o == "" {
		o = "date desc"
	}
	result := database.DB
	if query.Title != "" {
		// if title is given, do case insensitive search in title string
		result = result.Where("position(LOWER(?) in LOWER(title))>0", query.Title)
	} else {
		// else search exact query match
		result = result.Where(query)
	}

	if !date_from.IsZero() {
		result = result.Where("date > ?", date_from)
	}

	if !date_to.IsZero() {
		result = result.Where("date < ?", date_to)
	}

	result = result.Order(o).Find(c).Scopes(Paginate(c, len(*c), &pagination, database.DB)).Preload("Creator").Preload("Game").Preload("Vod").Find(c)
	if result.RowsAffected == 0 {
		return &pagination, errors.New("not found")
	}
	return &pagination, nil
}

func AddNewClip(c *models.Clip) error {
	if c.UUID == "" {
		b := make([]byte, 3)
		for {
			rand.Read(b)
			uuid := hex.EncodeToString(b)
			if err := GetOneClip(&models.Clip{}, uuid); err != nil {
				c.UUID = uuid
				break
			}
		}
	}

	var creator models.Creator
	var game models.Game
	var vod models.Vod
	omits := []string{}

	if err := GetOneCreator(&creator, c.CreatorUUID); err != nil {
		omits = append(omits, "CreatorUUID")
	}
	if err := GetOneGame(&game, c.GameUUID); err != nil {
		omits = append(omits, "GameUUID")
	}
	if err := GetOneVod(&vod, c.VodUUID, false); err != nil {
		omits = append(omits, "VodUUID")
	}
	if err := database.DB.Omit(strings.Join(omits, ",")).Create(&c).Error; err != nil {
		return err
	}
	return nil
}

func GetOneClip(c *models.Clip, uuid string) error {
	result := database.DB.Where("uuid = ?", uuid).Preload("Creator").Preload("Game").Preload("Vod").Find(c)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetClipByFilename(c *models.Clip, filename string) error {
	result := database.DB.Where("filename = ?", filename).Preload("Creator").Preload("Game").Preload("Vod").Find(c)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetClipsByUUID(c *[]models.Clip, uuids []string) error {
	result := database.DB.Where("uuid IN ?", uuids).Preload("Creator").Preload("Game").Preload("Vod").Find(c)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func PatchClip(changes map[string]interface{}, uuid string) error {
	var clip models.Clip
	if err := GetOneClip(&clip, uuid); err != nil {
		return errors.New("clip not found")
	}
	if err := database.DB.Model(&clip).Where("uuid = ?", uuid).Updates(changes).Error; err != nil {
		return errors.New("update failed")
	}
	return nil
}

func DeleteClip(c *models.Clip, uuid string) error {
	database.DB.Where("uuid = ?", uuid).Delete(c)
	return nil
}
