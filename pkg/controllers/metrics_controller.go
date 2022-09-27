package controllers

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

// GetMetrics godoc
// @Summary Get process infos
// @Tags    Metrics
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router  /metrics/ [get]
func GetMetrics(c *gin.Context) {
	var metrics runtime.MemStats
	runtime.ReadMemStats(&metrics)

	c.JSON(http.StatusOK, gin.H{
		"CpuCores":   runtime.NumCPU(),
		"alloc":      metrics.Alloc,
		"TotalAlloc": metrics.TotalAlloc,
		"Sys":        metrics.Sys,
		"NumGC":      metrics.NumGC,
		"GoVersion":  runtime.Version(),
	})
}
