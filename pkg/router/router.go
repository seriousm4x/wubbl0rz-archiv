package router

import (
	"os"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/logger"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"
var AuthMiddleware *jwt.GinJWTMiddleware

type User struct {
	Name string
}

func Init() *gin.Engine {
	// Create auth middleware
	var err error
	jwt_secret_key_timeout, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_TIMEOUT"))
	AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "kekl zone",
		Key:         []byte(os.Getenv("JWT_SECRET_KEY")),
		Timeout:     time.Duration(jwt_secret_key_timeout) * time.Hour,
		MaxRefresh:  time.Duration(jwt_secret_key_timeout) * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Name: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			if username == os.Getenv("JWT_USER") && password == os.Getenv("JWT_PASSWORD") {
				return &User{
					Name: username,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.Name == os.Getenv("JWT_USER") {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		logger.Error.Fatal("[router] JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := AuthMiddleware.MiddlewareInit()

	if errInit != nil {
		logger.Error.Fatal("[router] authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	// register routes
	r := gin.New()
	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/health"),
		gin.Recovery(),
	)
	r.Static("/media", "/var/www/media")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := r.Group("/")
	PublicRoutes(v1)
	PrivateRoutes(v1)

	return r
}
