package controllers

import (
	"net/http"
	"strconv"

	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"
	"github.com/gin-gonic/gin"
)

// Get Emotes godoc
// @Summary Get all emotes
// @Tags    Emotes
// @Accept  json
// @Produce  json
// @Success 200 {array}  models.Emote
// @Failure  400 {string} string
// @Failure  404 {string} string
// @Router  /emotes/ [get]
// @Param   id query string    false "The id of an emote"
// @Param   name query string false "The name of an emote"
// @Param   provider query string false "The provider of an emote"
func GetEmotes(c *gin.Context) {
	var emotes []models.Emote
	var query models.Emote

	c.ShouldBindQuery(&query)

	// set pagination
	var pagination queries.Pagination
	limit, _ := strconv.Atoi(c.Query("limit"))
	pagination.Limit = limit
	page, _ := strconv.Atoi(c.Query("page"))
	pagination.Page = page

	page_obj, err := queries.GetAllEmotes(&emotes, query, pagination)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error":  true,
			"msg":    "No emotes found",
			"result": emotes,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"msg":      "Ok",
		"page_obj": page_obj,
		"result":   emotes,
	})
}

// CreateEmote godoc
// @Summary  Create emote
// @Tags     Emotes
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  201 {string} string
// @Failure  400 {string} string
// @Failure  422 {string} string
// @Router   /emotes/ [post]
// @Param    Body body models.Emote true "Emote obj"
func CreateEmote(c *gin.Context) {
	var newEmote models.Emote
	var emote models.Emote

	if err := c.BindJSON(&newEmote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	if err := queries.GetOneEmote(&emote, newEmote.ID); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Emote already exists",
		})
		return
	}

	if err := queries.AddNewEmote(&newEmote); err != nil {
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

// PatchEmote godoc
// @Summary  Patch emote
// @Tags     Emotes
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure 400 {string} string
// @Failure  422 {string} string
// @Router   /emotes/{id} [patch]
// @Param    id path int         true "Unique Identifier"
// @Param    Body body models.Emote true "Emote obj"
func PatchEmote(c *gin.Context) {
	id := c.Param("id")

	// use map[string]interface{} for body values, so that gorm updates zero values
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var params map[string]interface{}
	c.ShouldBind(&params)

	if err := queries.PatchEmote(params, id); err != nil {
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

// DeleteEmote godoc
// @Summary  Delete emote
// @Tags     Emotes
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router   /emotes/{id} [delete]
// @Param    id path string true "Unique Identifier"
func DeleteEmote(c *gin.Context) {
	var emote models.Emote
	id := c.Param("id")

	if err := queries.GetOneEmote(&emote, id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": true,
			"msg":   "Emote not found",
		})
		return
	}

	if err := queries.DeleteEmote(&emote, id); err != nil {
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
