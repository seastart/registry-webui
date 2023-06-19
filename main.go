package main

import (
	"embed"
	"flag"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/seastart/registry-webui/lib"
)

//go:embed ui/dist
var ui embed.FS

func main() {
	// parse command line arguments
	var (
		confFilePath = ""
	)
	// for development, please copy default.yml to local.yml and run `go run main.go`
	flag.StringVar(&confFilePath, "config", "./config/local.yml", "config file path")
	flag.Parse()

	config, logger := lib.Init(confFilePath)
	// time ticker to fresh all repoes
	ticker := time.NewTicker(config.GetDuration("app.fresh_interval"))
	go func() {
		// init all repoes on start
		lib.Repoes.Refresh()
		for range ticker.C {
			lib.Repoes.Refresh()
		}
	}()
	// gin
	if config.GetString("app.env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(ginzap.Ginzap(logger.Desugar(), time.RFC3339, false))
	r.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))
	r.Use(cors.Default())
	// all route to vue
	// like nginx try_files
	r.NoRoute(func(c *gin.Context) {
		prefix := "ui/dist"
		path := c.Request.URL.Path // req path
		ext := filepath.Ext(path)  // file extension
		// read file data
		if data, err := ui.ReadFile(prefix + path); err != nil {
			// if file not exists, return index.html
			if data, err = ui.ReadFile(prefix + "/index.html"); err != nil {
				c.JSON(404, gin.H{
					"err": err,
				})
			} else {
				c.Data(200, mime.TypeByExtension(".html"), data)
			}
		} else {
			// if file exists, set mime type and return file content
			c.Data(200, mime.TypeByExtension(ext), data)
		}
	})
	// app info
	r.GET("/api/app", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"title": config.GetString("app.title"),
		})
	})
	// refresh all repoes
	r.POST("/api/refresh", func(c *gin.Context) {
		err := lib.Repoes.Refresh()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data": "refresh success",
			})
		}
	})
	// repo list
	r.GET("/api/repoes", func(c *gin.Context) {
		// keyword
		keyword := c.Query("keyword")
		// pagination
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(c.Query("per-page"))
		if err != nil || pageSize <= 0 || pageSize > 100 {
			pageSize = 100
		}
		repoes, hasmore, err := lib.Repoes.GetPage(keyword, page, pageSize)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":     repoes,
				"has_more": hasmore,
			})
		}
	})
	// repo detail
	r.GET("/api/repo", func(c *gin.Context) {
		name := c.Query("name")
		refresh := c.Query("refresh")
		repo, err := lib.Repoes.GetDetail(name, refresh == "1")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data": repo,
			})
		}
	})
	// pprof
	pprof.Register(r)
	r.Run(":8081")
}
