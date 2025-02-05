package utils

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func ConvertMarkdownToGoogleDocs(filePath string, removeMd bool, errors chan error) {
    newFilePath := strings.Replace(filePath, ".md", ".docx", 1)
    cmd := exec.Command("pandoc", "-f", "markdown", "-t", "docx", filePath, "-o", newFilePath)
    cmdOutput, err := cmd.CombinedOutput()
    if err != nil {
        errors <- fmt.Errorf("pandoc failed: %v\nOutput: %s", err, string(cmdOutput))
        return
    }

    if removeMd {
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            errors <- fmt.Errorf("file does not exist: %s", filePath)
            return
        }

        if err := os.Remove(filePath); err != nil {
            errors <- err
            return
        }
    }

    fmt.Println("Converted", filePath, "to", newFilePath)
}

func ConvertMarkdownInDirectoryToGoogleDocs(directoryPath string, removeMd bool) {
	LogInfo("Converting markdown files in " + directoryPath)
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(files))
	pool := make(chan struct{}, runtime.NumCPU())
	for _, file := range files {
		if file.IsDir() {
			ConvertMarkdownInDirectoryToGoogleDocs(directoryPath+"/"+file.Name(), removeMd)
			continue
		}
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()
			defer func() {
				<-pool
			}()
			pool <- struct{}{}

			fmt.Println(file.Name(), directoryPath+"/"+file.Name())
			content, err := os.ReadFile(directoryPath + "/" + file.Name())
			if err != nil {
				fmt.Println(err)
				errors <- err
				return
			}
			// Replace the "uploads/*"" in file content with the relative path
			newContent := strings.ReplaceAll(string(content), "uploads/", escapePath("./"+directoryPath+"/uploads/"))
			if err := os.WriteFile(directoryPath+"/"+file.Name(), []byte(newContent), 0644); err != nil {
				fmt.Println(err)
				errors <- err
				return
			}

			ConvertMarkdownToGoogleDocs(directoryPath+"/"+file.Name(), removeMd, errors)
		}(file)
	}

	wg.Wait()
	close(errors)
	for err := range errors {
		LogError(fmt.Sprint(err))
	}
}

func escapePath(path string) string {
	return url.PathEscape(path)
}
