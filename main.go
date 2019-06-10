package main

import (
	"fileserver/api"
	"fileserver/config"
	"fileserver/fs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	configuration := config.New()

	router := gin.Default()
	fs.Init(configuration.Root)

	router.Use(cors.Default())
	router.Use(gin.Logger())

	api.NewRouter(router)

	router.Run()
}
