package api

import (
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewRouter(router *gin.Engine) {
	h := &Handler{}

	v1 := router.Group("/v1")
	{
		v1.POST("/login", h.login)

		v1.GET("/*path", h.getList)
		v1.POST("/*path", h.upload)
		v1.DELETE("/*path", h.delete)
	}
}

func abortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
