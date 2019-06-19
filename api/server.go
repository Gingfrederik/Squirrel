package api

import (
	"fileserver/auth"
	authMiddleware "fileserver/middleware/auth"
	"fileserver/middleware/jwt"
	"fileserver/types"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewRouter(router *gin.Engine) {
	authM := auth.GetInstance()

	h := &Handler{}

	v1 := router.Group("/v1")
	{
		uR := v1.Group("/user")
		{
			uR.POST("/login", h.login)
			uR.POST("/register", h.register)
			uR.GET("/list", jwt.JWTAuth(), authMiddleware.NewAuthorizer(authM.Enforcer), h.getAllUser)
		}

		fsR := v1.Group("/fs")
		fsR.Use(jwt.JWTAuth())
		fsR.Use(authMiddleware.NewAuthorizer(authM.Enforcer))
		{
			fsR.GET("/*path", h.getList)
			fsR.POST("/*path", h.upload)
			fsR.DELETE("/*path", h.delete)
		}

		acR := v1.Group("/ac")
		acR.Use(jwt.JWTAuth())
		acR.Use(authMiddleware.NewAuthorizer(authM.Enforcer))
		{
			acR.GET("/policy", h.getAllPolicy)
			acR.POST("/policy", h.addPolicy)
			acR.DELETE("/policy", h.delPolicy)

			acR.GET("/role", h.getAllRole)

			acR.GET("/role/user", h.getAllUserRole)
			acR.POST("/role/user", h.addUserRole)
			acR.DELETE("/role/user", h.delUserRole)
		}
	}
}

func abortWithError(c *gin.Context, code int, message string) {
	res := types.Response{
		Status:  -1,
		Message: message,
	}
	c.AbortWithStatusJSON(code, res)
}
