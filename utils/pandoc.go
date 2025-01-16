package utils

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func ConvertMarkdownToGoogleDocs(filePath string, removeMd bool, errors chan error) {
	newFilePath := strings.Replace(filePath, ".md", ".docx", 1)
	cmd := exec.Command("pandoc", "-f", "markdown", "-t", "docx", filePath, "-o", newFilePath)
	if err := cmd.Run(); err != nil {
		errors <- err
		return
	}

	if removeMd {
		if err := os.Remove(filePath); err != nil {
			errors <- err
			return
		}
	}

	fmt.Println("Converted", filePath, "to", newFilePath)
}

func ConvertMarkdownInDirectoryToGoogleDocs(directoryPath string, removeMd bool) {
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	errors := make(chan error)
	pool := make(chan struct{}, 5)
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
			fmt.Println(file.Name(), directoryPath+"/"+file.Name())
			pool <- struct{}{}

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
		fmt.Println(err)
	}
}

func escapePath(path string) string {
	return url.PathEscape(path)
}
