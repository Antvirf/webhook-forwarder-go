package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

// Run HTTP server on specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Create new server instance + create routes
func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)
	server := &Server{}
	router := gin.Default()

	// aliveness check
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "alive!\n")
	})

	// take a webhook, validate, and forward
	router.POST("/forward_webhook", server.forwardWebhook)

	// print out received data
	router.POST("/receive_webhook", func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)
		ctx.JSON(http.StatusOK, string(body))
		//log.Print(string(body))
		//log.Print(ctx.Request.Header)
	})

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
