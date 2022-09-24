package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func createSignature(token string, message string) string {
	signature := hmac.New(sha256.New, []byte(token))
	signature.Write([]byte(message))
	return hex.EncodeToString(signature.Sum(nil))
}

func AddHeadersToRequest(ctx *gin.Context, req *http.Request) {
	for key, header := range ctx.Request.Header {
		// if v is list, add multiple times
		for _, value := range header {
			req.Header.Add(key, value)
		}
	}
}

func getIpRangesFromApiResponse(response map[string]interface{}) []net.IPNet {
	var ips []net.IPNet
	for k, v := range response {
		if k == "hooks" {
			for _, value := range v.([]interface{}) {
				_, adderValue, _ := net.ParseCIDR(value.(string))
				ips = append(ips, *adderValue)
			}
		}
	}
	return ips
}

func CheckIpInAcceptedRange(incomingIps []string, allowedIps []net.IPNet) bool {
	for _, strIncIp := range incomingIps {
		incIp := net.ParseIP(strIncIp)
		for _, allowedIpRange := range allowedIps {
			if allowedIpRange.Contains(incIp) {
				return true
			}
		}
	}
	return false
}

func GetGithubMetaApiFromFile(filepath string) map[string]interface{} {
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}

// file finders from https://stackoverflow.com/questions/60584697/find-a-file-by-regex-in-golang-given-the-regex-and-path
func walkFiles(root string, fn func(string) bool) []string {
	var files []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if fn(s) {
			files = append(files, s)
		}
		return nil
	})
	return files
}

// file finders from https://stackoverflow.com/questions/60584697/find-a-file-by-regex-in-golang-given-the-regex-and-path
func FindAndCleanJsonFiles() (string, error) {
	var returnFiles []string
	files := walkFiles(".", func(s string) bool {
		return filepath.Ext(s) == ".json"
	})
	for _, file := range files {
		// get date
		splitFilename := strings.Split(file, "_")
		splitFilenameStr := splitFilename[len(splitFilename)-1]
		splitFilenameStr = strings.Split(splitFilenameStr, ".")[0]

		// convert date
		date, err := time.Parse("20060102", splitFilenameStr)
		if err != nil {
			log.Fatal(err)
		}

		// check if older than 24 hours
		isExpired := date.Before(time.Now().Add(time.Hour * -24))

		// delete if yes
		if isExpired {
			os.Remove(file)
		} else {
			returnFiles = append(returnFiles, file)
		}
	}

	if len(returnFiles) >= 1 {
		returnFile := returnFiles[len(returnFiles)-1]
		log.Println("reading github api response from local file:", returnFile)
		return returnFile, nil
	} else {
		log.Println("no valid local api response json found, querying github...")
		return "", errors.New("no valid files found")
	}
}

// adapted from https://stackoverflow.com/questions/16311232/how-to-pipe-an-http-response-to-a-file-in-go
func queryGithubMetaApiToJson() string {
	resp, err := http.Get("https://api.github.com/meta")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	filepath := fmt.Sprintf("./github_meta_api_response_%s.json", time.Now().Format("20060102"))
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal("failed to create file")
	}
	defer out.Close()
	io.Copy(out, resp.Body)
	return filepath
}
