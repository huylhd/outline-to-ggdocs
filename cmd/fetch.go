package cmd

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
	"outline-to-ggdocs/config"
	"io/ioutil"
	"strings"
)

func getFileOperations(page int) ([]interface{}, error) {
	url := "https://app.getoutline.com/api/fileOperations.list"
	payloadString := `{"limit": 50, "offset": ` + fmt.Sprint((page-1)*50) + `, "sort": "updatedAt", "direction": "DESC", "type": "export"}`
	requestBody := strings.NewReader(payloadString)

	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+config.AppConfig.OutlineApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Debug: Print the request details
	fmt.Printf("Request URL: %s\n", req.URL)
	fmt.Printf("Request Headers: %v\n", req.Header)
	fmt.Printf("Request Body: %s\n", payloadString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	// Read the response body for debugging
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(responseBody, &bodyMap)
	if err != nil {
		return nil, err
	}

	data, ok := bodyMap["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	return data, nil
}


func removeDuplicates(fileOperations []interface{}) []interface{} {
    seen := make(map[string]bool)
    uniqueFileOperations := []interface{}{}

	for _, fo := range fileOperations {
		foMap, ok := fo.(map[string]interface{})
		if !ok {
			continue
		}

		collectionID, ok := foMap["collectionId"].(string)
		if !ok || collectionID == "" {
			continue
		}

		state, ok := foMap["state"].(string)
		if !ok || state == "expired" {
			continue
		}

		if !seen[collectionID] {
			seen[collectionID] = true
			uniqueFileOperations = append(uniqueFileOperations, fo)
		}
	}

    return uniqueFileOperations
}

func saveToFile(fileOperations []interface{}) error {
    file, err := os.Create("exportedFiles.json")
    if err != nil {
        return err
    }
    defer file.Close()

    var fileOperationData []map[string]interface{}
    for _, fo := range fileOperations {
        foMap := fo.(map[string]interface{})
        fileOperationData = append(fileOperationData, map[string]interface{}{
            "id":           foMap["id"],
            "name":         foMap["name"],
            "collectionId": foMap["collectionId"],
            "size":         foMap["size"],
        })
    }

    data, err := json.MarshalIndent(fileOperationData, "", "  ")
    if err != nil {
        return err
    }

    _, err = file.Write(data)
    return err
}

func FetchFileOperations() {
	var allFileOperations []interface{}
	page := 1

	for {
		fileOperations, err := getFileOperations(page)
		if err != nil {
			fmt.Printf("Error fetching file operations: %v\n", err)
			return
		}

		if len(fileOperations) == 0 {
			break
		}

		allFileOperations = append(allFileOperations, fileOperations...)
		page++
	}

	uniqueFileOperations := removeDuplicates(allFileOperations)

    err := saveToFile(uniqueFileOperations)
    if err != nil {
        fmt.Printf("Error saving file operations: %v\n", err)
        return
    }

    fmt.Println("File operations saved successfully to exportedFiles.json")
}