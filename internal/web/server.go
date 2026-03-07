package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"cheaptrick/internal/store"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

// NewRouter creates and configures the Gin routing engine
func NewRouter(reqStore *store.Store, fixtureDir string) *gin.Engine {
	r := gin.Default()

	hub := NewHub()
	go hub.Run()
	reqStore.Register(hub)

	h := &apiHandler{
		reqStore:   reqStore,
		fixtureDir: fixtureDir,
	}

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		api.GET("/requests", h.listRequests)
		api.GET("/requests/:id", h.getRequest)
		api.POST("/requests/:id/respond", h.respondToRequest)
		api.POST("/requests/:id/fixture", h.saveFixture)
		api.GET("/templates", h.listTemplates)
		api.GET("/fixtures", h.listFixtures)
		api.GET("/fixtures/:hash", h.getFixture)
		api.POST("/fixtures", h.createFixture)
		api.DELETE("/fixtures/:hash", h.deleteFixture)
		api.DELETE("/requests/:id", h.deleteRequest)
		api.DELETE("/requests", h.clearRequests)
	}

	// WebSocket
	r.GET("/ws", hub.HandleWebSocket)

	// Serve React SPA — MUST be last
	if os.Getenv("CHEAPTRICK_DEV") != "1" {
		distFS, err := fs.Sub(frontendFS, "frontend/dist")
		if err == nil {
			r.NoRoute(gin.WrapH(http.FileServer(http.FS(distFS))))
		}
	} else {
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		})
	}

	return r
}
