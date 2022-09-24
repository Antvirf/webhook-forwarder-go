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

func create_signature(token string, message string) string {
	signature := hmac.New(sha256.New, []byte(token))
	signature.Write([]byte(message))
	return hex.EncodeToString(signature.Sum(nil))
}

func add_headers_to_request(ctx *gin.Context, req *http.Request) {
	for key, header := range ctx.Request.Header {
		// if v is list, add multiple times
		for _, value := range header {
			req.Header.Add(key, value)
		}
	}
}

func get_ip_ranges_from_api_response(response map[string]interface{}) []net.IPNet {
	var ips []net.IPNet
	for k, v := range response {
		if k == "hooks" {
			for _, value := range v.([]interface{}) {
				_, adder_value, _ := net.ParseCIDR(value.(string))
				ips = append(ips, *adder_value)
			}
		}
	}
	return ips
}

func check_ip_in_accepted_range(incomingIps []string, allowedIps []net.IPNet) bool {
	for _, str_inc_ip := range incomingIps {
		inc_ip := net.ParseIP(str_inc_ip)
		for _, allowed_ip_range := range allowedIps {
			if allowed_ip_range.Contains(inc_ip) {
				// log.Println("SUCCESS:    ", inc_ip, "not in", allowed_ip_range)
				return true
			} else {
				// log.Println("FAIL:    ", inc_ip, "not in", allowed_ip_range)
			}
		}
	}
	return false
}

func get_github_meta_api_from_file(filepath string) map[string]interface{} {
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
func walk_files(root string, fn func(string) bool) []string {
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
func find_and_clean_json_files() (string, error) {
	var returning_files []string
	files := walk_files(".", func(s string) bool {
		return filepath.Ext(s) == ".json"
	})
	for _, file := range files {
		// get date
		split_filename := strings.Split(file, "_")
		split_filename_str := split_filename[len(split_filename)-1]
		split_filename_str = strings.Split(split_filename_str, ".")[0]

		// convert date
		date, err := time.Parse("20060102", split_filename_str)
		if err != nil {
			log.Fatal(err)
		}

		// check if older than 24 hours
		is_expired := date.Before(time.Now().Add(time.Hour * -24))

		// delete if yes
		if is_expired {
			os.Remove(file)
		} else {
			returning_files = append(returning_files, file)
		}
	}

	if len(returning_files) >= 1 {
		return_file := returning_files[len(returning_files)-1]
		log.Println("reading github api response from local file:", return_file)
		return return_file, nil
	} else {
		log.Println("no valid local api response json found, querying github...")
		return "", errors.New("no valid files found")
	}
}

// adapted from https://stackoverflow.com/questions/16311232/how-to-pipe-an-http-response-to-a-file-in-go
func query_github_meta_api_to_json() string {
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
