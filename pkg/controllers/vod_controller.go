package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"
	"github.com/gin-gonic/gin"
)

// GetVods godoc
// @Summary Get all vods
// @Tags    Vods
// @Accept  json
// @Produce  json
// @Success 200 {array}  models.Vod
// @Failure  400 {string} string
// @Failure  404 {string} string
// @Router  /vods/ [get]
// @Param   uuid       query string false "The uuid of a vod"
// @Param   title      query string false "The title of a vod"
// @Param   duration   query int    false "The duration of a vod"
// @Param   date       query string false "The date of a vod"
// @Param   filename   query string false "The filename of a vod"
// @Param   resolution query string false "The resolution of a vod"
// @Param   fps        query int    false "The fps of a vod"
// @Param   size       query int    false "The size of a vod"
// @Param   order      query string false "Set order direction divided by comma. Possible ordering values: 'date', 'duration', 'size'. Possible directions: 'asc', 'desc'. Example: 'date,desc'"
func GetVods(c *gin.Context) {
	var vods []models.Vod
	var query models.Vod

	c.ShouldBindQuery(&query)

	// filter array of uuid's
	// /vods/?uuids=aaa,bbb,ccc
	if c.Query("uuids") != "" {
		uuids := strings.Split(c.Query("uuids"), ",")
		err := queries.GetVodsByUUID(&vods, uuids)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":  true,
				"msg":    "uuid not found",
				"result": vods,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error":  false,
			"msg":    "Ok",
			"result": vods,
		})
		return
	}

	// /vods?year=2020
	if c.Query("year") != "" {
		year := c.Query("year")
		err := queries.GetVodsByYear(&vods, year)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":  true,
				"msg":    "no vods in year found",
				"result": vods,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error":  false,
			"msg":    "Ok",
			"result": vods,
		})
		return
	}

	// custom ordering query
	orderParams := ""
	if orderParams = c.Query("order"); orderParams != "" {
		order := strings.Split(orderParams, ",")
		if len(order) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "Invalid order params. Example: 'date,desc'",
			})
			return
		}
		if !stringInSlice(order[0], []string{"date", "duration", "size", "viewcount"}) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "Invalid first order param. 'date', 'duration' or 'size'",
			})
			return
		}
		if !stringInSlice(order[1], []string{"asc", "desc"}) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "Invalid second order param. 'asc' or 'desc'",
			})
			return
		}
		orderParams = strings.Replace(orderParams, ",", " ", -1)
		orderParams = orderParams + " nulls last"
	}

	// set pagination
	var pagination queries.Pagination
	limit, _ := strconv.Atoi(c.Query("limit"))
	pagination.Limit = limit
	page, _ := strconv.Atoi(c.Query("page"))
	pagination.Page = page

	page_obj, err := queries.GetAllVods(&vods, query, pagination, orderParams)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  true,
			"msg":    "No vods found",
			"result": vods,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"msg":      "Ok",
		"page_obj": page_obj,
		"result":   vods,
	})
}

// GetVodByUUID godoc
// @Summary Get vod by uuid
// @Tags    Vods
// @Produce  json
// @Success 200 {object} models.Vod
// @Failure 404 {string} string
// @Router  /vods/{uuid} [get]
// @Param   uuid path string true "Unique Identifier"
func GetVodByUUID(c *gin.Context) {
	var vod models.Vod

	if err := queries.GetOneVod(&vod, c.Param("uuid"), true); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  true,
			"msg":    "Vod not found",
			"result": vod,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":  false,
		"msg":    "Ok",
		"result": vod,
	})
}

// CreateVod godoc
// @Summary  Create vod
// @Tags     Vods
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  201 {string} string
// @Failure  400 {string} string
// @Failure  422 {string} string
// @Router   /vods/ [post]
// @Param    Body body models.Vod true "Vod obj"
func CreateVod(c *gin.Context) {
	var newVod models.Vod
	var vod models.Vod

	if err := c.BindJSON(&newVod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	if err := queries.GetOneVod(&vod, newVod.UUID, false); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Vod already exists",
		})
		return
	}

	if err := queries.AddNewVod(&newVod); err != nil {
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

// PatchVod godoc
// @Summary  Patch vod
// @Tags     Vods
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure  400 {string} string
// @Failure  422 {string} string
// @Router   /vods/{uuid} [patch]
// @Param    uuid path string     true "Unique Identifier"
// @Param    Body body models.Vod true "Vod obj"
func PatchVod(c *gin.Context) {
	// use map[string]interface{} for body values, so that gorm updates zero values
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var params map[string]interface{}
	c.ShouldBind(&params)
	uuid := c.Param("uuid")

	if err := queries.PatchVod(params, uuid); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": true,
			"msg":   "Failed to patch model",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msg":   "Updated",
	})
}

// DeleteVod godoc
// @Summary  Delete vod
// @Tags     Vods
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router   /vods/{uuid} [delete]
// @Param    uuid path string true "Unique Identifier"
func DeleteVod(c *gin.Context) {
	var vod models.Vod
	uuid := c.Param("uuid")

	if err := queries.GetOneVod(&vod, uuid, false); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,
			"msg":   "Vod not found",
		})
		return
	}

	if err := queries.DeleteVod(&vod, uuid); err != nil {
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
