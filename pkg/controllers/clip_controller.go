package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/AgileProggers/archiv-backend-go/pkg/queries"
	"github.com/gin-gonic/gin"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetClips godoc
// @Summary Get all clips
// @Tags    Clips
// @Accept  json
// @Produce  json
// @Success 200 {array}  models.Clip
// @Failure  400 {string} string
// @Failure  404 {string} string
// @Router  /clips/ [get]
// @Param   uuid       query string false "The uuid of a clip"
// @Param   title      query string false "The title of a clip"
// @Param   duration   query int    false "The duration of a clip"
// @Param   date       query string false "The date of a clip"
// @Param   filename   query string false "The filename of a clip"
// @Param   resolution query string false "The resolution of a clip"
// @Param   size       query int    false "The size of a clip"
// @Param   viewcount  query int    false "The viewcount of a clip"
// @Param   creator    query int    false "The creator id of a clip"
// @Param   game       query int    false "The game id of a clip"
// @Param   vod        query string false "The vod id of a clip"
// @Param   order      query string false "Set order direction divided by comma. Possible ordering values: 'date', 'duration', 'size'. Possible directions: 'asc', 'desc'. Example: 'date,desc'"
func GetClips(c *gin.Context) {
	var clips []models.Clip
	var query models.Clip

	c.ShouldBindQuery(&query)

	// filter array of uuid's
	// /clips/?uuids=aaa,bbb,ccc
	if c.Query("uuids") != "" {
		uuids := strings.Split(c.Query("uuids"), ",")
		if err := queries.GetClipsByUUID(&clips, uuids); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  true,
				"msg":    "Clip not found",
				"result": clips,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error":  false,
			"msg":    "Ok",
			"result": clips,
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
				"msg":   "Invalid first order param. 'date', 'duration', 'size' or 'viewcount'",
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
	}

	// custom date range
	var date_from time.Time
	var date_to time.Time

	query_date_from := c.Query("date_from")
	if query_date_from != "" {
		parsed_date, err := time.Parse(time.RFC3339, query_date_from)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "'date_from' not in RFC3339 format",
			})
			return
		}
		date_from = parsed_date
	}
	query_date_to := c.Query("date_to")
	if query_date_to != "" {
		parsed_date, err := time.Parse(time.RFC3339, query_date_to)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "'date_to' not in RFC3339 format",
			})
			return
		}
		date_to = parsed_date
	}

	// set pagination
	var pagination queries.Pagination
	limit, _ := strconv.Atoi(c.Query("limit"))
	pagination.Limit = limit
	page, _ := strconv.Atoi(c.Query("page"))
	pagination.Page = page

	page_obj, err := queries.GetAllClips(&clips, query, pagination, orderParams, date_from, date_to)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  true,
			"msg":    "No clips found",
			"result": clips,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"msg":      "Ok",
		"page_obj": page_obj,
		"result":   clips,
	})
}

// GetClipByID godoc
// @Summary Get clips by uuid
// @Tags    Clips
// @Produce  json
// @Success 200 {object} models.Clip
// @Failure 404 {string} string
// @Router  /clips/{uuid} [get]
// @Param   uuid path string true "Unique Identifier"
func GetClipByUUID(c *gin.Context) {
	var clip models.Clip

	if err := queries.GetOneClip(&clip, c.Param("uuid")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  true,
			"msg":    "Clip not found",
			"result": clip,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":  false,
		"msg":    "Ok",
		"result": clip,
	})
}

// CreateClip godoc
// @Summary  Create clip
// @Tags     Clips
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  201 {string} string
// @Failure  400 {string} string
// @Failure  422 {string} string
// @Router   /clips/ [post]
// @Param    Body body models.Clip true "Clip obj"
func CreateClip(c *gin.Context) {
	var newClip models.Clip
	var clip models.Clip

	if err := c.BindJSON(&newClip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	if err := queries.GetOneClip(&clip, newClip.UUID); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Clip already exists",
		})
		return
	}

	if err := queries.AddNewClip(&newClip); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error": false,
		"msg":   "Created",
	})
}

// PatchClip godoc
// @Summary  Patch clip
// @Tags     Clips
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure  400 {string} string
// @Failure  422 {string} string
// @Router   /clips/{uuid} [patch]
// @Param    uuid path string      true "Unique Identifier"
// @Param    Body body models.Clip true "Clip obj"
func PatchClip(c *gin.Context) {
	// use map[string]interface{} for body values, so that gorm updates zero values
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var params map[string]interface{}
	c.ShouldBind(&params)
	uuid := c.Param("uuid")

	if err := queries.PatchClip(params, uuid); err != nil {
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

// DeleteClip godoc
// @Summary  Delete clip
// @Tags     Clips
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Success  200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router   /clips/{uuid} [delete]
// @Param    uuid path string true "Unique Identifier"
func DeleteClip(c *gin.Context) {
	var clip models.Clip
	uuid := c.Param("uuid")

	if err := queries.GetOneClip(&clip, uuid); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,
			"msg":   "Clip not found",
		})
		return
	}

	if err := queries.DeleteClip(&clip, uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Error while deleting the model",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msg":   "Deleted",
	})
}
