package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"outline-to-ggdocs/config"
	"outline-to-ggdocs/constants"
	"outline-to-ggdocs/utils"
)

func ExportCollection(id string) map[string]interface{} {
	if id == "" {
		utils.LogInfo("Exporting all collections...")
	} else {
		utils.LogInfo(fmt.Sprint("Exporting collection with ID:", id))
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

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		utils.LogError("failed to export collections.")
		utils.LogError("Header: " + fmt.Sprint(resp.Header))
		utils.LogError("Body: " + fmt.Sprint(bodyMap))
		panic("failed to export collections")
	}

	data := bodyMap["data"].(map[string]interface{})
	return data
}

func ExportCollectionsCommand(id string) {
	data := ExportCollection(id)
	fileOperationId := data["fileOperation"].(map[string]interface{})["id"].(string)
	fmt.Println("Export requested, File operation ID:", fileOperationId)
	if utils.ShouldProceedInput("Download file? (y/n) ") {
		DownloadFile(fileOperationId, "./")
	}
}
