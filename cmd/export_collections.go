package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"outline-to-ggdocs/config"
	"outline-to-ggdocs/constants"
)

func ExportCollectionsCommand(id string) {
	if id == "" {
		fmt.Println("Exporting all collections...")
	} else {
		fmt.Println("Exporting collection with ID:", id)
	}

	apiUrl := constants.OutlineApiCollectionsExportAll

	payload := make(map[string]string)
	payload["format"] = "outline-markdown"
	if id != "" {
		apiUrl = constants.OutlineApiCollectionsExport
		payload["id"] = id
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonPayload))
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

	var bodyMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&bodyMap)

	if resp.StatusCode != 200 {
		panic("Failed to export collections")
	}

	if err != nil {
		panic(err)
	}

	data := bodyMap["data"].(map[string]interface{})
	fileOperationId := data["fileOperation"].(map[string]interface{})["id"].(string)
	fmt.Println("Export requested, File operation ID:", fileOperationId)
	fmt.Print("Download file? (y/n) ")
	var download string
	fmt.Scanln(&download)
	if download != "y" {
		return
	}
	DownloadFile(fileOperationId)
}
