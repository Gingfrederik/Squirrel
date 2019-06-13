package api

import (
	"fileserver/middleware/jwt"
	"fileserver/types"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	genID *snowflake.Node
}

func NewRouter(router *gin.Engine) {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	h := &Handler{
		genID: node,
	}

	v1 := router.Group("/v1")
	{
		uR := v1.Group("/user")
		{
			uR.POST("/login", h.login)
			uR.POST("/register", h.register)
		}

		fsR := v1.Group("/fs")
		fsR.Use(jwt.JWTAuth())
		{
			fsR.GET("/*path", h.getList)
			fsR.POST("/*path", h.upload)
			fsR.DELETE("/*path", h.delete)
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
