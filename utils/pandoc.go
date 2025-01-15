package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func ConvertMarkdownToGoogleDocs(filePath string, removeMd bool, wg *sync.WaitGroup, errors chan error) {
	defer wg.Done()

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
			pool <- struct{}{}
			ConvertMarkdownToGoogleDocs(directoryPath+"/"+file.Name(), removeMd, &wg, errors)
			<-pool
		}(file)
	}

	wg.Wait()
	close(errors)
	for err := range errors {
		fmt.Println(err)
	}
}
