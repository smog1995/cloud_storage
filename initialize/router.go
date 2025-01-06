package initialize

import (
	"cloud_storage/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Routers() *gin.Engine {
	zap.S().Debugf("初始化路由器")
	Router := gin.Default()
	fileGroup := Router.Group("/file")
	{
		fileGroup.GET("/upload", handler.UploadHandler)
		fileGroup.POST("/upload", handler.UploadHandler)
		fileGroup.POST("/meta", handler.GetFileMetaHandler)
		fileGroup.POST("update", hander.Upda)
	}
	userGroup := Router.Group("/user")
	{
		userGroup.GET("/signup", handler.SignupHandler)
		userGroup.POST("/signup", handler.SignupHandler)
		userGroup.GET("/login", handler.LoginHandler)
		userGroup.POST("/login", handler.LoginHandler)
		userGroup.POST("/userinfo", handler.UserInfoHandler)
	}

	Router.Use(Cors())
	return Router
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
