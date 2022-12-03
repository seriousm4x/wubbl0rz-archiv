package router

import (
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/controllers"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

var MemoryStore *persist.MemoryStore

// PublicRoutes func for describe group of public routes.
func PublicRoutes(rg *gin.RouterGroup) {
	// caching
	MemoryStore = persist.NewMemoryStore(24 * time.Hour)

	// auth login
	rg.POST("/token/new", AuthMiddleware.LoginHandler)

	// other routes
	rg.GET("/download/:type/:uuid", controllers.SendStream)
	rg.GET("/health", controllers.GetHealth)
	rg.GET("/chat", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetAllChatMessages)

	// route groups
	vodsGroup := rg.Group("/vods")
	{
		vodsGroup.GET("/", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetVods)
		vodsGroup.GET("/:uuid", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetVodByUUID)
	}
	clipsGroup := rg.Group("/clips")
	{
		clipsGroup.GET("/", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetClips)
		clipsGroup.GET("/:uuid", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetClipByUUID)
	}
	gamesGroup := rg.Group("/games")
	{
		gamesGroup.GET("/", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetGames)
		gamesGroup.GET("/:uuid", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetGameByUUID)
	}
	creatorsGroup := rg.Group("/creators")
	{
		creatorsGroup.GET("/", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetCreators)
		creatorsGroup.GET("/:uuid", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetCreatorByUUID)
	}
	emotesGroup := rg.Group("/emotes")
	{
		emotesGroup.GET("/", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetEmotes)
	}
	yearsGroup := rg.Group("/years")
	{
		yearsGroup.GET("/", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetYears)
	}
	statsGroup := rg.Group("/stats")
	{
		statsGroup.GET("/short", controllers.GetShortStats)
		statsGroup.GET("/long", cache.CacheByRequestURI(MemoryStore, 24*time.Hour), controllers.GetLongStats)
	}
}
