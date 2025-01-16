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
	"outline-to-ggdocs/utils"
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
		panic("failed to fetch file info")
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

func DownloadFile(id string, rootDir string) string {
	if id == "" {
		panic("id is required")
	}

	filePath := rootDir
	for {
		fmt.Printf("Retrieving file %s status...\n", id)
		ready, collectionName, err := isFileReady(id)
		if err != nil {
			utils.LogError(fmt.Sprintf("Error retrieving file %s status: %s", id, err))
			break
		}
		if !ready {
			fmt.Println("File is not ready yet, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		payloadString := `{"id": "` + id + `"}`
		req, err := http.NewRequest("POST", constants.OutlineApiFileOperationsDownload, strings.NewReader(payloadString))
		if err != nil {
			utils.LogError(fmt.Sprintf("Error creating request for file %s: %s", id, err))
			break
		}

		req.Header.Add("Authorization", "Bearer "+config.AppConfig.OutlineApiKey)
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{
			Timeout: 300 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			utils.LogError(fmt.Sprintf("Error downloading file %s: %s", id, err))
			break
		}
		defer resp.Body.Close()

		filePath += "/" + fmt.Sprintf("%s.zip", collectionName)
		if resp.StatusCode == 200 {
			fmt.Printf("Downloading file %s \n", id)
			file, err := os.Create(filePath)
			if err != nil {
				utils.LogError(fmt.Sprintf("Error creating file %s: %s", filePath, err))
				break
			}
			defer file.Close()

			buf := make([]byte, 1024*32)
			_, err = io.CopyBuffer(file, resp.Body, buf)
			if err != nil {
				utils.LogError(fmt.Sprintf("Error writing file %s: %s", filePath, err))
				break
			}

			utils.LogInfo(fmt.Sprintf("Download file %s completed \n", id))
			break
		}

		utils.LogError(fmt.Sprintf("Download file %s failed \n", id))
		break
	}

	return filePath
}
