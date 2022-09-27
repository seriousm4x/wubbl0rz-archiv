package controllers

import (
	"net/http"
	"strconv"

	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"
	"github.com/gin-gonic/gin"
)

// GetCreators godoc
// @Summary Get all creators
// @Tags    Creators
// @Accept  json
// @Produce  json
// @Success 200 {array}  models.Creator
// @Failure  400 {string} string
// @Failure  404 {string} string
// @Router  /creators/ [get]
// @Param   uuid query int    false "The uuid of a creator"
// @Param   name query string false "The name of a creator"
func GetCreators(c *gin.Context) {
	var creators []models.Creator
	var query models.Creator

	c.ShouldBindQuery(&query)

	// set pagination
	var pagination queries.Pagination
	limit, _ := strconv.Atoi(c.Query("limit"))
	pagination.Limit = limit
	page, _ := strconv.Atoi(c.Query("page"))
	pagination.Page = page

	page_obj, err := queries.GetAllCreators(&creators, query, pagination)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  true,
			"msg":    "No creators found",
			"result": creators,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"msg":      "Ok",
		"page_obj": page_obj,
		"result":   creators,
	})
}

// GetCreatorByID godoc
// @Summary Get creator by uuid
// @Tags    Creators
// @Produce  json
// @Success 200 {object} models.Creator
// @Failure  400 {string} string
// @Failure 404 {string} string
// @Router  /creators/{uuid} [get]
// @Param   uuid path int true "Unique Identifyer"
func GetCreatorByUUID(c *gin.Context) {
	var creator models.Creator
	uuid := c.Param("uuid")

	if err := queries.GetOneCreator(&creator, uuid); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  true,
			"msg":    "Creator not found",
			"result": creator,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":  false,
		"msg":    "Ok",
		"result": creator,
	})
}

// CreateCreator godoc
// @Summary  Create creator
// @Tags     Creators
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  201 {string} string
// @Failure  400 {string} string
// @Failure  422 {string} string
// @Router   /creators/ [post]
// @Param    Body body models.Creator true "Creator obj"
func CreateCreator(c *gin.Context) {
	var newCreator models.Creator
	var creator models.Creator

	if err := c.BindJSON(&newCreator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	if err := queries.GetOneCreator(&creator, newCreator.UUID); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Creator already exists",
		})
		return
	}

	if err := queries.AddNewCreator(&newCreator); err != nil {
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

// PatchCreator godoc
// @Summary  Patch creator
// @Tags     Creators
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure 400 {string} string
// @Failure  422 {string} string
// @Router   /creators/{uuid} [patch]
// @Param    uuid path int            true "Unique Identifier"
// @Param    Body body models.Creator true "Creator obj"
func PatchCreator(c *gin.Context) {
	uuid := c.Param("uuid")

	// use map[string]interface{} for body values, so that gorm updates zero values
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var params map[string]interface{}
	c.ShouldBind(&params)

	if err := queries.PatchCreator(params, uuid); err != nil {
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

// DeleteCreator godoc
// @Summary  Delete creator
// @Tags     Creators
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router   /creators/{uuid} [delete]
// @Param    uuid path string true "Unique Identifier"
func DeleteCreator(c *gin.Context) {
	var creator models.Creator
	uuid := c.Param("uuid")

	if err := queries.GetOneCreator(&creator, uuid); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,
			"msg":   "Creator not found",
		})
		return
	}

	if err := queries.DeleteCreator(&creator, uuid); err != nil {
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
