package queries

import (
	"encoding/json"
	"errors"

	"github.com/AgileProggers/archiv-backend-go/pkg/database"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
)

func GetAllEmotes(e *[]models.Emote, query models.Emote, pagination Pagination) (*Pagination, error) {
	result := database.DB
	if query.Name != "" {
		// if name is given, do case insensitive search in name string
		result = result.Where("position(LOWER(?) in LOWER(name))>0", query.Name)
	} else {
		// else search exact query match
		result = result.Where(query)
	}

	result = result.Find(e).Scopes(Paginate(e, len(*e), &pagination, database.DB)).Find(e)
	if result.RowsAffected == 0 {
		return &pagination, errors.New("not found")
	}
	return &pagination, nil
}

func AddNewEmote(e *models.Emote) error {
	if err := database.DB.Create(e).Error; err != nil {
		return err
	}
	return nil
}

func GetOneEmote(e *models.Emote, id string) error {
	result := database.DB.Where("id = ?", id).Find(e)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func PatchEmote(changes map[string]interface{}, id string) error {
	var emote models.Emote
	if err := GetOneEmote(&emote, id); err != nil {
		return errors.New("emote not found")
	}
	if err := database.DB.Model(&emote).Where("id = ?", id).Updates(changes).Error; err != nil {
		return errors.New("update failed")
	}
	return nil
}

func DeleteEmote(e *models.Emote, id string) error {
	database.DB.Where("id = ?", id).Delete(e)
	return nil
}

func UpdateOrCreateEmote(e *models.Emote, id string) error {
	var changes map[string]interface{}
	data, _ := json.Marshal(e)
	json.Unmarshal(data, &changes)
	changes["outdated"] = false

	result := database.DB.Model(&e).Where("id = ?", id).Where("provider = ?", e.Provider).Updates(&changes)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		if err := AddNewEmote(e); err != nil {
			return err
		}
	}
	return nil
}

func MarkAllEmotesOutdated(outdated bool) error {
	var emotes models.Emote
	if result := database.DB.Model(&emotes).Where("outdated != ?", outdated).Updates(map[string]interface{}{"outdated": outdated}); result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteOutdatedEmotes() error {
	var emotes models.Emote
	if result := database.DB.Model(&emotes).Where("outdated = ?", true).Delete(&emotes); result.Error != nil {
		return result.Error
	}
	return nil
}
