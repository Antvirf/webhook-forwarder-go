package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func getGithubMetaApiResponse() map[string]interface{} {
	// step 1: loop through json, delete if older than 1 day
	// if returns err, then need to download. otherwise, returns latest
	filepath, err := FindAndCleanJsonFiles()

	// step 2: if it returned an error, we need to fetch it from GitHub
	if err != nil {
		filepath = queryGithubMetaApiToJson()
	}

	// step 3: read file and return map
	return GetGithubMetaApiFromFile(filepath)

}

func (server *Server) forwardWebhook(ctx *gin.Context) {

	apiResponse := getGithubMetaApiResponse()
	hookIpRanges := getIpRangesFromApiResponse(apiResponse)

	// define list of incoming ips
	incomingIps := []string{
		ctx.Request.Host,
		ctx.Request.Header.Get("Origin"),
		ctx.Request.Header.Get("x-client-ip"),
		ctx.Request.Header.Get("x-forwarded-for"),
	}

	// check one: x-client-ip / x-forwarded-for ip
	conditionIpMatch := CheckIpInAcceptedRange(incomingIps, hookIpRanges)
	if !conditionIpMatch {
		log.Println("unauthorized sender IPs", incomingIps, "returning 401")
		ctx.JSON(http.StatusUnauthorized, "unauthorized IP")
		return
	}

	// check three: webhook signature
	requestBody, _ := io.ReadAll(ctx.Request.Body)
	hookSignature := ctx.Request.Header.Get("x-hub-signature-256")

	secret_value := os.Getenv("WEBHOOK_TOKEN_SECRET")
	computedSignature := createSignature(secret_value, string(requestBody))

	// compare digest hmac compare digest local sig vs payload sig
	if hookSignature != computedSignature {
		log.Println("signature does not match, returning 401")
		ctx.JSON(http.StatusUnauthorized, "unauthorized webhook")
		return
	}

	// create new request for forwarding
	TARGET_URL := os.Getenv("TARGET_URL")
	req, err := http.NewRequest("POST", TARGET_URL, bytes.NewBuffer(requestBody))
	if err != nil {
		errorResponse(err)
	}

	// add headers from context to the new request
	AddHeadersToRequest(ctx, req)

	// create client and forward the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorResponse(err)
	}
	defer resp.Body.Close()

	ctx.JSON(http.StatusOK, "hook forwarded!")
}
