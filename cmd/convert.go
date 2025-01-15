package cmd

import "outline-to-ggdocs/utils"

func ConvertDirectory(directoryPath string, removeMd bool) {
	utils.ConvertMarkdownInDirectoryToGoogleDocs(directoryPath, removeMd)
}
