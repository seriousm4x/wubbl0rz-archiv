package queries

import (
	"encoding/json"
	"errors"

	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/database"
	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/models"
)

func GetSettings(s *models.Settings) error {
	s.ID = 1
	result := database.DB.Model(&s).Where("id = ?", s.ID).Find(&s)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no settings")
	}
	return nil
}

func InitSettings(s *models.Settings) error {
	s.ID = 1
	if result := database.DB.Model(&s).Create(&s); result.Error != nil {
		return result.Error
	}
	return nil
}

func PartiallyUpdateSettings(s *models.Settings) error {
	s.ID = 1
	if result := database.DB.Model(&s).Where("id = ?", s.ID).Updates(&s); result.Error != nil {
		return result.Error
	}
	return nil
}

func OverwriteAllSettings(s *models.Settings) error {
	s.ID = 1
	// count settings rows
	var count int64
	result := database.DB.Model(&s).Where("id = ?", s.ID).Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		// update
		var changes map[string]interface{}
		data, _ := json.Marshal(s)
		json.Unmarshal(data, &changes)
		if resUpdate := result.Updates(&changes); resUpdate.Error != nil {
			return resUpdate.Error
		}
	} else {
		// create
		if result := database.DB.Model(&s).Create(&s); result.Error != nil {
			return result.Error
		}
	}
	return nil
}
