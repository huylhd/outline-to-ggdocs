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

func AutomateCommand(page int, fromStep AutomateStep, toStep AutomateStep) {
	if fromStep <= STEP_LIST && toStep >= STEP_LIST {
		stepList(page)
	}
	if fromStep <= STEP_EXPORT && toStep >= STEP_EXPORT {
		stepExport()
	}
	if fromStep <= STEP_DOWNLOAD_UNZIP_CONVERT && toStep >= STEP_DOWNLOAD_UNZIP_CONVERT {
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
	collectionsData := readCollectionsData()
	if len(collectionsData) == 0 {
		utils.LogError("no collections found in collections.json")
		os.Exit(0)
	}

	exportedFiles := readDataFromFile("exportedFiles.json")
	exportedCollectionIds := make(map[string]bool)
	for _, file := range exportedFiles {
		exportedCollectionIds[file["collectionId"].(string)] = true
	}

	var wg sync.WaitGroup
	ch := make(chan map[string]interface{}, len(collectionsData))
	for _, c := range collectionsData {
		if exportedCollectionIds[c["id"].(string)] {
			continue
		}
		wg.Add(1)
		go func(c map[string]interface{}) {
			defer wg.Done()
			data := ExportCollection(c["id"].(string))
			fileOperationId := data["fileOperation"].(map[string]interface{})["id"].(string)
			fileOperationData := map[string]interface{}{
				"id":           fileOperationId,
				"name":         c["name"].(string),
				"collectionId": c["id"].(string),
				"size":         data["fileOperation"].(map[string]interface{})["size"].(string),
			}
			ch <- fileOperationData

			exportedFile := map[string]interface{}{
				"collectionId": c["id"].(string),
				"id":           fileOperationId,
				"name":         c["name"].(string),
				"size":         data["fileOperation"].(map[string]interface{})["size"].(string),
			}
			appendDataFile("exportedFiles.json", exportedFile)
		}(c)
	}

	wg.Wait()
	close(ch)
	fileOperationData := make([]map[string]interface{}, 0)
	for i := range ch {
		fileOperationData = append(fileOperationData, i)
	}
	file, err := os.Create("fileOperations.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(fileOperationData, "", "  ")
	if err != nil {
		panic(err)
	}
	file.Write(jsonData)

	utils.LogInfo("Collections export requested successfully")
}



func appendDataFile(fileName string, downloadedFile map[string]interface{}) {
    err := utils.AppendDataToFile(fileName, downloadedFile)
    if err != nil {
        panic(err)
    }
}

func readDataFromFile(fileName string) []map[string]interface{} {
    data, err := utils.ReadDataFromFile(fileName)
    if err != nil {
        panic(err)
    }
    return data
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
	file, err := os.Open("fileOperations.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var fileOperationData []map[string]interface{}
	err = json.NewDecoder(file).Decode(&fileOperationData)
	if err != nil {
		panic(err)
	}

	targetPath := "exports"

	os.Mkdir("exports", 0755)
	var wg sync.WaitGroup
	for _, i := range fileOperationData {
		wg.Add(1)
		fileId := i["id"].(string)
		collectionId := i["collectionId"].(string)
		name := i["name"].(string)
		go func(fileId string, collectionId string, name string) {
			defer wg.Done()
			filePath, err := DownloadFile(fileId, "exports")
			if err != nil {
				return
			}
			utils.Unzip(filePath, targetPath)
			
			utils.LogInfo("Downloaded and unzipped collection: " + name)
			downloadedFile := map[string]interface{}{
				"fileId":       fileId,
				"collectionId": collectionId,
				"name":		 name,
			}
			appendDataFile("downloadedFiles.json", downloadedFile)
		
		}(fileId, collectionId, name)
	}

	wg.Wait()

	// convert markdown files to google docs
	utils.ConvertMarkdownInDirectoryToGoogleDocs(targetPath, true)
	utils.LogInfo("All collections downloaded and converted successfully")
}
