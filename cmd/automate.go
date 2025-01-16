package cmd

import (
	"encoding/json"
	"fmt"
	"outline-to-ggdocs/utils"
	"sync"

	"os"
)

type AutomateStep int

const (
	STEP_LIST AutomateStep = iota
	STEP_EXPORT
	STEP_DOWNLOAD_UNZIP_CONVERT
)

func AutomateCommand(page int, fromStep AutomateStep) {
	if fromStep <= STEP_LIST {
		stepList(page)
	}
	if fromStep <= STEP_EXPORT {
		stepExport()
	}
	if fromStep <= STEP_DOWNLOAD_UNZIP_CONVERT {
		stepDownloadAndUnzip()
	}
}

func stepList(page int) {
	utils.LogInfo("Fetching all collections")
	collections := listAllCollections(page)
	if len(collections) == 0 {
		utils.LogError("No collections found")
		os.Exit(0)
	}

	fmt.Println("Saving collection list to collections.json")
	storedCollections := make([]map[string]interface{}, 0)
	for _, collection := range collections {
		id := collection.(map[string]interface{})["id"].(string)
		name := collection.(map[string]interface{})["name"].(string)
		c := make(map[string]interface{})
		c["id"] = id
		c["name"] = name
		storedCollections = append(storedCollections, c)
	}
	file, err := os.Create("collections.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(storedCollections, "", "  ")
	if err != nil {
		panic(err)
	}
	file.Write(jsonData)

	fmt.Println("*You can modify collections.json to export specific collections")
	if !utils.ShouldProceedInput("Export collections? (y/n) ") {
		os.Exit(0)
	}
}

func listAllCollections(p int) []interface{} {
	var collections []interface{}

	page := p
	if page == 0 {
		page = 1
	}
	for {
		c := ListCollections(page)
		if len(c) == 0 {
			break
		}
		collections = append(collections, c...)
		if p == 0 {
			break
		}
		page++
	}

	return collections
}

func stepExport() {
	var collectionsData []map[string]interface{}
	collectionsData = readCollectionsData()
	if len(collectionsData) == 0 {
		utils.LogError("no collections found in collections.json")
		os.Exit(0)
	}
	var wg sync.WaitGroup
	ch := make(chan string, len(collectionsData))
	for _, c := range collectionsData {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := ExportCollection(c["id"].(string))
			fileOperationId := data["fileOperation"].(map[string]interface{})["id"].(string)
			ch <- fileOperationId
		}()
	}

	wg.Wait()
	close(ch)
	fileOperationIds := make([]string, 0)
	for fileOperationId := range ch {
		fileOperationIds = append(fileOperationIds, fileOperationId)
	}
	file, err := os.Create("fileOperationIds.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(fileOperationIds, "", "  ")
	if err != nil {
		panic(err)
	}
	file.Write(jsonData)

	utils.LogInfo("Collections export requested successfully")
}

func readCollectionsData() []map[string]interface{} {
	file, err := os.Open("collections.json")
	if err != nil {
		panic(err)
	}

	var collectionsData []map[string]interface{}
	err = json.NewDecoder(file).Decode(&collectionsData)
	if err != nil {
		panic(err)
	}

	return collectionsData
}

func stepDownloadAndUnzip() {
	file, err := os.Open("fileOperationIds.json")
	if err != nil {
		panic(err)
	}

	var fileOperationIds []string
	err = json.NewDecoder(file).Decode(&fileOperationIds)
	if err != nil {
		panic(err)
	}

	os.Mkdir("exports", 0755)
	var wg sync.WaitGroup
	for _, fileId := range fileOperationIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filePath := DownloadFile(fileId, "exports")
			targetPath := "exports"
			utils.Unzip(filePath, targetPath)
			utils.ConvertMarkdownInDirectoryToGoogleDocs(targetPath, true)
		}()
	}

	wg.Wait()
	utils.LogInfo("All collections downloaded and converted successfully")
}
