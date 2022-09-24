package api

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func get_github_meta_api_response() map[string]interface{} {
	// step 1: loop through jsons, delete if older than 1 day
	// if returns err, then need to download. otherwise, returns latest
	filepath, err := find_and_clean_json_files()

	// step 2: if it returned an error, we need to fetch it from GitHub
	if err != nil {
		filepath = query_github_meta_api_to_json()
	}

	// step 3: read file and return map
	return get_github_meta_api_from_file(filepath)

}

func (server *Server) forwardWebhook(ctx *gin.Context) {

	api_response := get_github_meta_api_response()
	hook_ip_ranges := get_ip_ranges_from_api_response(api_response)

	// define list of incoming ips
	incomingIps := []string{
		ctx.Request.Host,
		ctx.Request.Header.Get("Origin"),
		ctx.Request.Header.Get("x-client-ip"),
		ctx.Request.Header.Get("x-forwarded-for"),
	}

	// check one: x-client-ip / x-forwarded-for ip
	conditionIpMatch := check_ip_in_accepted_range(incomingIps, hook_ip_ranges)
	if !conditionIpMatch {
		log.Println("unauthorized sender IPs", incomingIps, "returning 401")
		ctx.JSON(http.StatusUnauthorized, "unauthorized IP")
		return
	}

	// check three: webhook signature
	requestBody, _ := io.ReadAll(ctx.Request.Body)
	hookSignature := ctx.Request.Header.Get("x-hub-signature-256")
	computedSignature := create_signature("secret", string(requestBody))

	// compare digest hmac compare digest local sig vs payload sig
	if !(hookSignature == computedSignature) {
		log.Println("signature does not match, returning 401")
		ctx.JSON(http.StatusUnauthorized, "unauthorized webhook")
		return
	}

	// create new request for forwarding
	req, err := http.NewRequest("POST", "http://localhost:8080/receive_webhook", bytes.NewBuffer(requestBody))
	if err != nil {
		errorResponse(err)
	}

	// add headers from context to the new request
	add_headers_to_request(ctx, req)

	// create client and forward the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorResponse(err)
	}
	defer resp.Body.Close()

	ctx.JSON(http.StatusOK, "hook forwarded!")
}
