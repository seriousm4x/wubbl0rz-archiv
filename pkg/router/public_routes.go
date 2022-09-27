package router

import (
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/controllers"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(rg *gin.RouterGroup) {
	// auth login
	rg.POST("/token/new", AuthMiddleware.LoginHandler)

	// other routes
	rg.GET("/download/:type/:uuid", controllers.SendStream)
	rg.GET("/health", controllers.GetHealth)

	// caching
	store := persistence.NewInMemoryStore(time.Second)

	// route groups
	vodsGroup := rg.Group("/vods")
	{
		vodsGroup.GET("/", cache.CachePage(store, 3*time.Hour, controllers.GetVods))
		vodsGroup.GET("/:uuid", cache.CachePage(store, 3*time.Hour, controllers.GetVodByUUID))
	}
	clipsGroup := rg.Group("/clips")
	{
		clipsGroup.GET("/", cache.CachePage(store, 3*time.Hour, controllers.GetClips))
		clipsGroup.GET("/:uuid", cache.CachePage(store, 3*time.Hour, controllers.GetClipByUUID))
	}
	gamesGroup := rg.Group("/games")
	{
		gamesGroup.GET("/", cache.CachePage(store, 3*time.Hour, controllers.GetGames))
		gamesGroup.GET("/:uuid", cache.CachePage(store, 3*time.Hour, controllers.GetGameByUUID))
	}
	creatorsGroup := rg.Group("/creators")
	{
		creatorsGroup.GET("/", cache.CachePage(store, 3*time.Hour, controllers.GetCreators))
		creatorsGroup.GET("/:uuid", cache.CachePage(store, 3*time.Hour, controllers.GetCreatorByUUID))
	}
	emotesGroup := rg.Group("/emotes")
	{
		emotesGroup.GET("/", cache.CachePage(store, 3*time.Hour, controllers.GetEmotes))
	}
	yearsGroup := rg.Group("/years")
	{
		yearsGroup.GET("/", cache.CachePage(store, 3*time.Hour, controllers.GetYears))
	}
	statsGroup := rg.Group("/stats")
	{
		statsGroup.GET("/short", controllers.GetShortStats)
		statsGroup.GET("/long", controllers.GetLongStats)
	}
}
