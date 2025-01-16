package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"outline-to-ggdocs/config"
	"outline-to-ggdocs/constants"
	"strings"
)

func ListCollections(page int) []interface{} {
	payloadString := `{"limit": 50, "offset": ` + fmt.Sprint((page-1)*100) + `}`
	req, err := http.NewRequest("POST", constants.OutlineApiCollectionsList, strings.NewReader(payloadString))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Failed to fetch collections %s \n", bodyMap))
	}

	collections := bodyMap["data"].([]interface{})
	return collections
}

func ListCollectionsCommand(page int) {
	collections := ListCollections(page)
	if len(collections) == 0 {
		fmt.Println("No collections found")
		return
	}
	for _, collection := range collections {
		collectionMap := collection.(map[string]interface{})

		fmt.Println("Collection ID:", collectionMap["id"], "| Name:", collectionMap["name"])
	}

	fmt.Print(">Enter the collection ID to export:")
	var id string
	fmt.Scanln(&id)
	ExportCollectionsCommand(id)
}
