package queries

import (
	"errors"

	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/database"
	"github.com/seriousm4x/wubbl0rz-archiv-transcribe/pkg/models"
)

func GetAllCreators(c *[]models.Creator, query models.Creator, pagination Pagination) (*Pagination, error) {
	result := database.DB.Model(&query)
	if query.Name != "" {
		// if name is given, do case insensitive search in name string
		result = result.Where("position(LOWER(?) in LOWER(name))>0", query.Name)
	} else {
		// else search exact query match
		result = result.Where(query)
	}
	result = result.Count(&pagination.TotalRows).Scopes(Paginate(&pagination, database.DB)).Find(c)
	if result.RowsAffected == 0 {
		return &pagination, errors.New("not found")
	}
	return &pagination, nil
}

func AddNewCreator(c *models.Creator) error {
	if err := database.DB.Create(c).Error; err != nil {
		return err
	}
	return nil
}

func GetOneCreator(c *models.Creator, uuid string) error {
	result := database.DB.Where("uuid = ?", uuid).Preload("Clips").Find(c)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func PatchCreator(changes map[string]interface{}, uuid string) error {
	var creator models.Creator
	if err := GetOneCreator(&creator, uuid); err != nil {
		return errors.New("creator not found")
	}
	if err := database.DB.Model(&creator).Where("uuid = ?", uuid).Updates(changes).Error; err != nil {
		return errors.New("update failed")
	}
	return nil
}

func DeleteCreator(c *models.Creator, uuid string) error {
	database.DB.Where("uuid = ?", uuid).Delete(c)
	return nil
}
