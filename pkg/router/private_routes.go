package router

import (
	"github.com/AgileProggers/archiv-backend-go/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func PrivateRoutes(rg *gin.RouterGroup) {
	rg.Use(AuthMiddleware.MiddlewareFunc())
	rg.GET("/token/refresh", AuthMiddleware.RefreshHandler)
	rg.GET("/metrics", controllers.GetMetrics)

	vodsGroup := rg.Group("/vods")
	{
		vodsGroup.POST("/", controllers.CreateVod)
		vodsGroup.PATCH("/:uuid", controllers.PatchVod)
		vodsGroup.DELETE("/:uuid", controllers.DeleteVod)
	}
	clipsGroup := rg.Group("/clips")
	{
		clipsGroup.POST("/", controllers.CreateClip)
		clipsGroup.PATCH("/:uuid", controllers.PatchClip)
		clipsGroup.DELETE("/:uuid", controllers.DeleteClip)
	}
	gamesGroup := rg.Group("/games")
	{
		gamesGroup.POST("/", controllers.CreateGame)
		gamesGroup.PATCH("/:uuid", controllers.PatchGame)
		gamesGroup.DELETE("/:uuid", controllers.DeleteGame)
	}
	creatorsGroup := rg.Group("/creators")
	{
		creatorsGroup.POST("/", controllers.CreateCreator)
		creatorsGroup.PATCH("/:uuid", controllers.PatchCreator)
		creatorsGroup.DELETE("/:uuid", controllers.DeleteCreator)
	}
	emotesGroup := rg.Group("/emotes")
	{
		emotesGroup.POST("/", controllers.CreateEmote)
		emotesGroup.PATCH("/:id", controllers.PatchEmote)
		emotesGroup.DELETE("/:id", controllers.DeleteEmote)
	}
}
