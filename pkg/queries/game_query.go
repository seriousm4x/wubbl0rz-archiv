package queries

import (
	"errors"

	"github.com/AgileProggers/archiv-backend-go/pkg/database"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
)

func GetAllGames(g *[]models.Game, query models.Game, pagination Pagination) (*Pagination, error) {
	result := database.DB
	if query.Name != "" {
		// if name is given, do case insensitive search in name string
		result = result.Where("position(LOWER(?) in LOWER(name))>0", query.Name)
	} else {
		// else search exact query match
		result = result.Where(query)
	}

	result = result.Find(g).Scopes(Paginate(g, len(*g), &pagination, database.DB)).Preload("Clips").Find(g)
	if result.RowsAffected == 0 {
		return &pagination, errors.New("not found")
	}
	return &pagination, nil
}

func AddNewGame(g *models.Game) error {
	if err := database.DB.Create(g).Error; err != nil {
		return err
	}
	return nil
}

func GetOneGame(g *models.Game, uuid string) error {
	result := database.DB.Where("uuid = ?", uuid).Preload("Clips").Find(g)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func PatchGame(changes map[string]interface{}, uuid string) error {
	var game models.Game
	if err := GetOneGame(&game, uuid); err != nil {
		return errors.New("game not found")
	}
	if err := database.DB.Model(&game).Where("uuid = ?", uuid).Updates(changes).Error; err != nil {
		return errors.New("update failed")
	}
	return nil
}

func DeleteGame(g *models.Game, uuid string) error {
	database.DB.Where("uuid = ?", uuid).Delete(g)
	return nil
}
