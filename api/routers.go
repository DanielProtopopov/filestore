package api

import (
	"filestore/api/group/public"
	"filestore/api/middleware"
	apiconfig "filestore/cmd/api/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func Serve() {
	route := gin.Default()

	route.LoadHTMLGlob("data/templates/*")

	route.Static("/fonts/Fira-Sans/", "data/fonts/Fira-Sans/")
	route.StaticFile("/css/style.css", "data/style.css")

	route.Use(middleware.NoCache).Use(LangWithConfig)
	route.GET("/", public.IndexHandler)
	route.GET("/api", public.APIHandler)
	route.GET("/about", public.AboutHandler)

	route.GET("/sitemap.xml", public.GetSitemapHandler)
	route.POST("/api/v1/upload", public.StoreTemporaryFileHandler)

	route.GET("/:path/:filename", public.GetTemporaryFileHandler)
	errRun := route.Run(fmt.Sprintf("%s:%d", apiconfig.Settings.Server.Host, apiconfig.Settings.Server.Port))
	log.Panic(errRun)
}

func LangWithConfig(c *gin.Context) {
	middleware.Language(c)
}
