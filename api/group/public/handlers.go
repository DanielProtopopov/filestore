package public

import (
	"context"
	apiconfig "filestore/cmd/api/config"
	"filestore/structs"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// GetSitemapHandler
// @Summary Get sitemap
// @Description Get sitemap
// @ID public-get-sitemap
// @Tags Public methods
// @Produce xml
// @Success 200 {} string string
// @Failure 400 {} string string
// @Failure 404 {} string string
// @Failure 500 {} string string
// @Router /sitemap/ [get]
func GetSitemapHandler(c *gin.Context) {
	rqContext := c.Request.Context()
	SentryHub := sentry.GetHubFromContext(rqContext)
	if SentryHub == nil {
		SentryHub = sentry.CurrentHub().Clone()
		rqContext = sentry.SetHubOnContext(rqContext, SentryHub)
	}

	Language, _ := c.Get("Language")
	Localizer := i18n.NewLocalizer(apiconfig.Settings.Bundle, Language.(string))

	rqContext = context.WithValue(rqContext, "Sentry", SentryHub)
	rqContext = context.WithValue(rqContext, "Localizer", Localizer)
	rqContext = context.WithValue(rqContext, "Language", Language)

	c.HTML(http.StatusOK, "", gin.H{})
}

// APIHandler
// @Summary Get api HTML
// @Description Get api HTML
// @ID public-get-api
// @Tags Public methods
// @Produce xml
// @Success 200 {} string string
// @Failure 400 {} string string
// @Failure 404 {} string string
// @Failure 500 {} string string
// @Router /api [get]
func APIHandler(c *gin.Context) {
	rqContext := c.Request.Context()
	SentryHub := sentry.GetHubFromContext(rqContext)
	if SentryHub == nil {
		SentryHub = sentry.CurrentHub().Clone()
		rqContext = sentry.SetHubOnContext(rqContext, SentryHub)
	}

	Language, _ := c.Get("Language")
	Localizer := i18n.NewLocalizer(apiconfig.Settings.Bundle, Language.(string))

	rqContext = context.WithValue(rqContext, "Sentry", SentryHub)
	rqContext = context.WithValue(rqContext, "Localizer", Localizer)
	rqContext = context.WithValue(rqContext, "Language", Language)

	c.HTML(http.StatusOK, "api.tmpl", gin.H{"uacode": apiconfig.Settings.GoogleTag, "protocol": apiconfig.Settings.Server.Protocol, "domain": apiconfig.Settings.Server.Domain,
		"maxuploadsize": apiconfig.Settings.MaximumUploadSize / (1024 * 1024), "filename": "", "folder": "", "expires": 0, "error": ""})
}

// AboutHandler
// @Summary Get about HTML
// @Description Get about HTML
// @ID public-get-about
// @Tags Public methods
// @Produce xml
// @Success 200 {} string string
// @Failure 400 {} string string
// @Failure 404 {} string string
// @Failure 500 {} string string
// @Router /about [get]
func AboutHandler(c *gin.Context) {
	rqContext := c.Request.Context()
	SentryHub := sentry.GetHubFromContext(rqContext)
	if SentryHub == nil {
		SentryHub = sentry.CurrentHub().Clone()
		rqContext = sentry.SetHubOnContext(rqContext, SentryHub)
	}

	Language, _ := c.Get("Language")
	Localizer := i18n.NewLocalizer(apiconfig.Settings.Bundle, Language.(string))

	rqContext = context.WithValue(rqContext, "Sentry", SentryHub)
	rqContext = context.WithValue(rqContext, "Localizer", Localizer)
	rqContext = context.WithValue(rqContext, "Language", Language)

	c.HTML(http.StatusOK, "about.tmpl", gin.H{"uacode": apiconfig.Settings.GoogleTag, "protocol": apiconfig.Settings.Server.Protocol, "domain": apiconfig.Settings.Server.Domain,
		"maxuploadsize": apiconfig.Settings.MaximumUploadSize / (1024 * 1024), "filename": "", "folder": "", "expires": 0, "error": ""})
}

// IndexHandler
// @Summary Get index HTML
// @Description Get index HTML
// @ID public-get-index
// @Tags Public methods
// @Produce xml
// @Success 200 {} string string
// @Failure 400 {} string string
// @Failure 404 {} string string
// @Failure 500 {} string string
// @Router / [get]
func IndexHandler(c *gin.Context) {
	rqContext := c.Request.Context()
	SentryHub := sentry.GetHubFromContext(rqContext)
	if SentryHub == nil {
		SentryHub = sentry.CurrentHub().Clone()
		rqContext = sentry.SetHubOnContext(rqContext, SentryHub)
	}

	Language, _ := c.Get("Language")
	Localizer := i18n.NewLocalizer(apiconfig.Settings.Bundle, Language.(string))

	rqContext = context.WithValue(rqContext, "Sentry", SentryHub)
	rqContext = context.WithValue(rqContext, "Localizer", Localizer)
	rqContext = context.WithValue(rqContext, "Language", Language)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"uacode": apiconfig.Settings.GoogleTag, "protocol": apiconfig.Settings.Server.Protocol, "domain": apiconfig.Settings.Server.Domain,
		"maxuploadsize": apiconfig.Settings.MaximumUploadSize / (1024 * 1024), "filename": "", "folder": "", "expires": 0, "error": ""})
}

// StoreTemporaryFileHandler
// @Summary Store temporary file and receive a link to it
// @Description Store temporary file and receive a link to it
// @ID public-store-temporary-file
// @Tags Public methods
// @Produce html
// @Success 200 {} string string
// @Failure 400 {} string string
// @Failure 404 {} string string
// @Failure 500 {} string string
// @Router /api/v1/upload [post]
func StoreTemporaryFileHandler(c *gin.Context) {
	rqContext := c.Request.Context()
	SentryHub := sentry.GetHubFromContext(rqContext)
	if SentryHub == nil {
		SentryHub = sentry.CurrentHub().Clone()
		rqContext = sentry.SetHubOnContext(rqContext, SentryHub)
	}

	Language, _ := c.Get("Language")
	Localizer := i18n.NewLocalizer(apiconfig.Settings.Bundle, Language.(string))

	rqContext = context.WithValue(rqContext, "Sentry", SentryHub)
	rqContext = context.WithValue(rqContext, "Localizer", Localizer)
	rqContext = context.WithValue(rqContext, "Language", Language)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, apiconfig.Settings.MaximumUploadSize)
	if errUploadFile := c.Request.ParseMultipartForm(apiconfig.Settings.MaximumUploadSize); errUploadFile != nil {
		if _, ok := errUploadFile.(*http.MaxBytesError); ok {
			c.JSON(http.StatusBadRequest, UploadRS{Status: fmt.Sprintf("File is too big: must be less than %d megabytes", apiconfig.Settings.MaximumExpirationPeriod), Errors: []string{fmt.Sprintf("File is too big: must be less than %d megabytes", apiconfig.Settings.MaximumExpirationPeriod/(1024*1024))}})
			return
		}

		c.JSON(http.StatusBadRequest, UploadRS{Status: fmt.Sprintf("Error uploading file: %s", errUploadFile.Error()), Errors: []string{fmt.Sprintf("Error uploading file: %s", errUploadFile.Error())}})
		return
	}

	storageDirectorySize, errGetStorageSize := DirectorySize(apiconfig.Settings.StoragePath)
	if errGetStorageSize != nil {
		log.Printf("[ERROR] Failed to get storage directory size: %s", errGetStorageSize.Error())
		c.JSON(http.StatusBadRequest, UploadRS{Status: fmt.Sprintf("Error getting storage directory size: %s", errGetStorageSize.Error()), Errors: []string{fmt.Sprintf("Error getting storage directory size: %s", errGetStorageSize.Error())}})
		return
	}
	if storageDirectorySize > apiconfig.Settings.MaximumStorageSize {
		log.Printf("[ERROR] Storage directory is full (%d bytes, limit %d bytes)", storageDirectorySize, apiconfig.Settings.MaximumStorageSize)
		c.JSON(http.StatusBadRequest, UploadRS{Status: fmt.Sprintf("Storage directory is full, please contact site administrator", apiconfig.Settings.MaximumStorageSize), Errors: []string{fmt.Sprintf("Storage directory is full, please contact site administrator", apiconfig.Settings.MaximumStorageSize)}})
		return
	}

	ipAddress := c.ClientIP()
	ipAddressFolders, errGetFolders := getIPAddressFolders(rqContext, apiconfig.Settings.StoragePath, ipAddress)
	if errGetFolders != nil {
		c.JSON(http.StatusBadRequest, UploadRS{Status: fmt.Sprintf("Error getting folders for IP address %s: %s", ipAddress, errGetFolders.Error()), Errors: []string{fmt.Sprintf("Error getting folders for IP address %s: %s", ipAddress, errGetFolders.Error())}})
		return
	}

	if uint(len(ipAddressFolders)) >= apiconfig.Settings.MaximumFilesPerIP {
		c.JSON(http.StatusBadRequest, UploadRS{Status: "You have uploaded the maximum number of files per IP address", Errors: []string{"You have uploaded the maximum number of files per IP address"}})
		return
	}

	expirationValue := c.DefaultPostForm("expire", "600")
	expirationSecondsVal, errConvertExpiration := strconv.Atoi(expirationValue)
	if errConvertExpiration != nil {
		c.JSON(http.StatusBadRequest, UploadRS{Status: fmt.Sprintf("Invalid expiry value: must be between 1 and %d seconds", apiconfig.Settings.MaximumExpirationPeriod),
			Errors: []string{fmt.Sprintf("Invalid expiry value: must be between 1 and %d seconds", apiconfig.Settings.MaximumExpirationPeriod)}})
		return
	}
	expirationSeconds := uint(expirationSecondsVal)
	if expirationSeconds < 1 {
		expirationSeconds = 60
	}
	if expirationSeconds >= apiconfig.Settings.MaximumExpirationPeriod {
		expirationSeconds = apiconfig.Settings.MaximumExpirationPeriod
	}
	_, fileHeader, errGetFile := c.Request.FormFile("file")
	if errGetFile != nil {
		c.JSON(http.StatusBadRequest, UploadRS{Status: "No file was supplied", Errors: []string{"No file was supplied"}})
		return
	}

	fileName := fileHeader.Filename
	randomFileFolders := generateFolderNames(10, 10)
	for _, randomFileFolder := range randomFileFolders {
		errCreateDirectory := os.Mkdir(filepath.Join(apiconfig.Settings.StoragePath, randomFileFolder), 0750)
		if errCreateDirectory != nil {
			// Done to avoid duplicates
			continue
		}

		errSaveFile := c.SaveUploadedFile(fileHeader, filepath.Join(apiconfig.Settings.StoragePath, randomFileFolder, fileHeader.Filename))
		if errSaveFile != nil {
			c.JSON(http.StatusBadRequest, UploadRS{Status: "Error saving the file, please try again",
				Errors: []string{fmt.Sprintf("Error saving the file: %s", errSaveFile.Error())}})
			return
		}

		// Save folder file
		structs.Redis.Set(rqContext, randomFileFolder, fileHeader.Filename, time.Duration(expirationSeconds)*time.Second)
		// Store IP address bound to folder with uploaded file for limit by IP address enforcement
		structs.Redis.LPush(rqContext, fmt.Sprintf("ip-%s", ipAddress), randomFileFolder)
		c.JSON(http.StatusOK, UploadRS{
			Errors: []string{}, Status: "success", Data: struct {
				Url string `json:"url"`
			}{Url: fmt.Sprintf("%s://%s/%s/%s", apiconfig.Settings.Server.Protocol, apiconfig.Settings.Server.Domain, randomFileFolder, fileName)},
		})
		return
	}

	c.JSON(http.StatusBadRequest, UploadRS{Status: "Failed to store the file, please try again", Errors: []string{"Failed to store the file, please try again"}})
}

// GetTemporaryFileHandler
// @Summary Retrieve a temporary file from a link
// @Description Retrieve a temporary file from a link
// @ID public-get-temporary-file
// @Tags Public methods
// @Produce html
// @Success 200 {} string string
// @Failure 400 {} string string
// @Failure 404 {} string string
// @Failure 500 {} string string
// @Router / [get]
func GetTemporaryFileHandler(c *gin.Context) {
	rqContext := c.Request.Context()
	SentryHub := sentry.GetHubFromContext(rqContext)
	if SentryHub == nil {
		SentryHub = sentry.CurrentHub().Clone()
		rqContext = sentry.SetHubOnContext(rqContext, SentryHub)
	}

	Language, _ := c.Get("Language")
	Localizer := i18n.NewLocalizer(apiconfig.Settings.Bundle, Language.(string))

	rqContext = context.WithValue(rqContext, "Sentry", SentryHub)
	rqContext = context.WithValue(rqContext, "Localizer", Localizer)
	rqContext = context.WithValue(rqContext, "Language", Language)

	path := c.Param("path")
	filename := c.Param("filename")
	if path == "" || filename == "" {
		c.HTML(http.StatusBadRequest, "index.tmpl", gin.H{"uacode": apiconfig.Settings.GoogleTag, "protocol": apiconfig.Settings.Server.Protocol, "domain": apiconfig.Settings.Server.Domain,
			"maxuploadsize": apiconfig.Settings.MaximumUploadSize / (1024 * 1024), "filename": "", "folder": "", "expires": 0, "error": ""})
		return
	}

	directories, errReadStoragePath := os.ReadDir(apiconfig.Settings.StoragePath)
	if errReadStoragePath == nil {
		for _, directory := range directories {
			if !directory.IsDir() {
				continue
			}
			if directory.Name() == path {
				c.FileAttachment(filepath.Join(apiconfig.Settings.StoragePath, path, filename), filename)
				return
			}
		}
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"uacode": apiconfig.Settings.GoogleTag, "protocol": apiconfig.Settings.Server.Protocol, "domain": apiconfig.Settings.Server.Domain,
		"maxuploadsize": apiconfig.Settings.MaximumUploadSize / (1024 * 1024), "filename": "", "folder": "", "expires": 0, "error": ""})
	return
}
