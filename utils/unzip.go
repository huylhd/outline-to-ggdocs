package utils

import (
	"fmt"
	"os/exec"
)

func Unzip(filePath string, targetPath string) {
	fmt.Println("Unzipping file", filePath)

	cmd := exec.Command("unzip", "-o", filePath, "-d", targetPath)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
