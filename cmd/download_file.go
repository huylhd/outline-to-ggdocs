package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"outline-to-ggdocs/config"
	"outline-to-ggdocs/constants"
	"strings"
	"time"
)

func isFileReady(id string) (bool, string, error) {
	payloadString := `{"id": "` + id + `"}`
	req, err := http.NewRequest("POST", constants.OutlineApiFileOperationsInfo, strings.NewReader(payloadString))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+config.AppConfig.OutlineApiKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Failed to fetch file info")
	}

	var bodyMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&bodyMap)
	if err != nil {
		panic(err)
	}

	data := bodyMap["data"].(map[string]interface{})
	state := data["state"].(string)
	if state == "creating" || state == "uploading" {
		return false, "", nil
	}
	if state == "complete" {
		return true, data["name"].(string), nil
	}
	return false, "", errors.New("get file operation failed")
}

func DownloadFile(id string) {
	if id == "" {
		panic("Error: id is required")
	}

	for {
		fmt.Println("Retrieving file status...")
		ready, collectionName, err := isFileReady(id)
		if err != nil {
			panic(err)
		}
		if !ready {
			fmt.Println("File is not ready yet, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		payloadString := `{"id": "` + id + `"}`
		req, err := http.NewRequest("POST", constants.OutlineApiFileOperationsDownload, strings.NewReader(payloadString))
		if err != nil {
			panic(err)
		}

		req.Header.Add("Authorization", "Bearer "+config.AppConfig.OutlineApiKey)
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			file, err := os.Create(fmt.Sprintf("%s.zip", collectionName))
			if err != nil {
				panic(err)
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				panic(err)
			}

			fmt.Println("Download completed")
			break
		}

		fmt.Println("Download failed")
		os.Exit(1)
	}
}
