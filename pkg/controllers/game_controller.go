package controllers

import (
	"net/http"
	"strconv"

	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"
	"github.com/gin-gonic/gin"
)

// Get Games godoc
// @Summary Get all games
// @Tags    Games
// @Accept  json
// @Produce  json
// @Success 200 {array}  models.Game
// @Failure  400 {string} string
// @Failure  404 {string} string
// @Router  /games/ [get]
// @Param   uuid    query int    false "The uuid of a game"
// @Param   name    query string false "The name of a game"
// @Param   box_art query string false "The box_art of a game"
func GetGames(c *gin.Context) {
	var games []models.Game
	var query models.Game

	c.ShouldBindQuery(&query)

	// set pagination
	var pagination queries.Pagination
	limit, _ := strconv.Atoi(c.Query("limit"))
	pagination.Limit = limit
	page, _ := strconv.Atoi(c.Query("page"))
	pagination.Page = page

	page_obj, err := queries.GetAllGames(&games, query, pagination)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error":  true,
			"msg":    "No games found",
			"result": games,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"msg":      "Ok",
		"page_obj": page_obj,
		"result":   games,
	})
}

// GetGameByID godoc
// @Summary Get game by uuid
// @Tags    Games
// @Produce  json
// @Success 200 {object} models.Game
// @Failure  400 {string} string
// @Failure 404 {string} string
// @Router  /games/{uuid} [get]
// @Param   uuid path string true "Unique Identifier"
func GetGameByUUID(c *gin.Context) {
	var game models.Game
	uuid := c.Param("uuid")

	if err := queries.GetOneGame(&game, uuid); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error":  true,
			"msg":    "Game not found",
			"result": game,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":  false,
		"msg":    "Ok",
		"result": game,
	})
}

// CreateGame godoc
// @Summary  Create game
// @Tags     Games
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  201 {string} string
// @Failure  400 {string} string
// @Failure  422 {string} string
// @Router   /games/ [post]
// @Param    Body body models.Game true "Game obj"
func CreateGame(c *gin.Context) {
	var newGame models.Game
	var game models.Game

	if err := c.BindJSON(&newGame); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	if err := queries.GetOneGame(&game, newGame.UUID); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Game already exists",
		})
		return
	}

	if err := queries.AddNewGame(&newGame); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": true,
			"msg":   "Error while creating the model",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error": false,
		"msg":   "Created",
	})
}

// PatchGame godoc
// @Summary  Patch game
// @Tags     Games
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure 400 {string} string
// @Failure  422 {string} string
// @Router   /games/{uuid} [patch]
// @Param    uuid path int         true "Unique Identifier"
// @Param    Body body models.Game true "Game obj"
func PatchGame(c *gin.Context) {
	uuid := c.Param("uuid")

	// use map[string]interface{} for body values, so that gorm updates zero values
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var params map[string]interface{}
	c.ShouldBind(&params)

	if err := queries.PatchGame(params, uuid); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": true,
			"msg":   "Error while patching the model",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msg":   "Updated",
	})
}

// DeleteGame godoc
// @Summary  Delete game
// @Tags     Games
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router   /games/{uuid} [delete]
// @Param    uuid path string true "Unique Identifier"
func DeleteGame(c *gin.Context) {
	var game models.Game
	uuid := c.Param("uuid")

	if err := queries.GetOneGame(&game, uuid); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": true,
			"msg":   "Game not found",
		})
		return
	}

	if err := queries.DeleteGame(&game, uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Error while deleting the model",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msg":   "Deleted",
	})
}
