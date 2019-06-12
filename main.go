package main

import (
	"fileserver/api"
	"fileserver/config"
	"fileserver/db"
	"fileserver/fs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.New()

	router := gin.Default()
	db.New(cfg.DB)
	fs.Init(cfg.Root)

	router.Use(cors.Default())
	router.Use(gin.Logger())

	api.NewRouter(router)

	router.Run()
}
