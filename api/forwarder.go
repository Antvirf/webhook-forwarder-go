package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// To be implemented for proper validation
type webhookBody struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

// To be implemented for proper validation
// TODO: Enforce IP structure for clean parsing
type webhookHeader struct {
	xClientIp     string `header:"X-Client-Ip"`
	xForwardedFor string `header:"x-forwarded-for"`
	xHubSignature string `header:"x-hub-signature-256"`
}

func (server *Server) forwardWebhook(ctx *gin.Context) {
	req := webhookBody{}
	// head := webhookHeader{}

	// // binding: HEADERS
	// if err := ctx.ShouldBindHeader(&head); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	return
	// }
	// binding: BODY
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get GitHub accepted IPs

	// check one: sender IP
	senderIp := ctx.Request.Header.Get("x-client-ip")
	fmt.Println("sender  ip:", senderIp)

	// check two: x-forwarded-for IP
	forwardedIp := ctx.Request.Header.Get("x-forwarded-for")
	fmt.Println("forwarded ip:", forwardedIp)

	// check three: webhook signature
	hookSignature := ctx.Request.Header.Get("x-hub-signature-256")
	fmt.Println("hook sig:", hookSignature)

	ctx.JSON(http.StatusOK, req)
}
