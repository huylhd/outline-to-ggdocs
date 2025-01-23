package utils

import (
    "encoding/json"
    "os"
    "fmt"
)

func ReadDataFromFile(fileName string) ([]map[string]interface{}, error) {
    // Open the file
	file, err := os.Open(fileName)
	if err != nil {
		// Handle the case where the file does not exist
		if os.IsNotExist(err) {
			file, err = os.Create(fileName)
			if err != nil {
				return nil, fmt.Errorf("failed to create file: %w", err)
			}
			defer file.Close()

		} else {
		    // Return other errors
		    return nil, fmt.Errorf("failed to open file: %w", err)
        }
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() == 0 {
		return []map[string]interface{}{}, nil
	}

    var data []map[string]interface{}
    err = json.NewDecoder(file).Decode(&data)
    if err != nil {
        return nil, err
    }

    return data, nil
}

func AppendDataToFile(fileName string, newData map[string]interface{}) error {
    data, err := ReadDataFromFile(fileName)
    if err != nil {
        return err
    }

    data = append(data, newData)

    file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    err = file.Truncate(0)
    if err != nil {
        return err
    }

    _, err = file.Seek(0, 0)
    if err != nil {
        return err
    }

    _, err = file.Write(jsonData)
    return err
}