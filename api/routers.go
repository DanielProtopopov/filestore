package api

import (
	"filestore/api/group/public"
	"filestore/api/middleware"
	apiconfig "filestore/cmd/api/config"
	"fmt"
	"log"
	"time"

	ratelimiter "github.com/khaaleoo/gin-rate-limiter/core"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Serve() {
	route := gin.Default()

	if len(apiconfig.Settings.TrustedProxies) > 0 {
		errSetTrustedProxies := route.SetTrustedProxies(apiconfig.Settings.TrustedProxies)
		if errSetTrustedProxies != nil {
			log.Panic("Failed to set trusted proxies list: %s", errSetTrustedProxies.Error())
		}
	}
	// max of 100 requests and then five more requests per second
	// the rate limiter will reset after 1 minute don't receive any requests by the same IP address
	rateLimiterOption := ratelimiter.RateLimiterOption{Limit: 5, Burst: 100, Len: 1 * time.Minute}

	// Create an IP rate limiter instance
	rateLimiterMiddleware := ratelimiter.RequireRateLimiter(ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter, Key: "iplimiter_maximum_requests_for_ip", Option: rateLimiterOption,
	})
	route.Use(gzip.Gzip(gzip.DefaultCompression)).Use(rateLimiterMiddleware)

	route.LoadHTMLGlob("data/templates/*")

	route.Static("/fonts/Fira-Sans/", "data/fonts/Fira-Sans/")
	route.StaticFile("/css/style.css", "data/style.css")
	route.StaticFile("/sitemap.xml", "data/web/sitemap.xml")
	route.StaticFile("/robots.txt", "data/web/robots.txt")
	route.StaticFile("/favicon.ico", "data/web/favicon.ico")

	route.Use(middleware.NoCache).Use(LangWithConfig)
	route.GET("/", public.IndexHandler)
	route.GET("/api", public.APIHandler)
	route.GET("/about", public.AboutHandler)

	route.POST("/api/v1/upload", public.StoreTemporaryFileHandler)
	route.GET("/:path/:filename", public.GetTemporaryFileHandler)

	errRun := route.Run(fmt.Sprintf("%s:%d", apiconfig.Settings.Server.Host, apiconfig.Settings.Server.Port))
	log.Panic(errRun)
}

func LangWithConfig(c *gin.Context) {
	middleware.Language(c)
}
