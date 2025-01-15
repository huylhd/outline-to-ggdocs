package main

import (
	"flag"
	"fmt"
	"os"
	"outline-to-ggdocs/cmd"
	"outline-to-ggdocs/config"
)

func main() {
	fmt.Println("--- Outline to Google Docs ---")

	flag.Usage = func() {
		fmt.Println("Usage: outline-to-ggdocs <command>")
		fmt.Println("Commands:")
		fmt.Println("  list - List all collections")
		fmt.Println("  export - Export a collection (or all without id)")
		fmt.Println("  download - Download a file")
		fmt.Println("  convert - Convert all markdown files in a directory to .docx")
		fmt.Println("\nRun 'outline-to-ggdocs <command> --help' to see more details about a specific command.")
	}

	args := os.Args[1:]
	help := flag.Bool("help", false, "Show help")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if len(args) < 1 {
		fmt.Println("Error: command is required")
		flag.Usage()
		os.Exit(1)
	}

	outlineApiKey := config.AppConfig.OutlineApiKey

	if outlineApiKey == "" {
		fmt.Println("Error: OUTLINE_API_KEY is required")
		os.Exit(1)
	}

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)

	switch args[0] {
	case "list":
		page := listCmd.Int("page", 1, "Page number")
		listCmd.Parse(args[1:])

		cmd.ListCollectionsCommand(*page)

	case "export":
		exportID := exportCmd.String("id", "", "Collection ID")
		exportCmd.Parse(args[1:])

		cmd.ExportCollectionsCommand(*exportID)

	case "download":
		downloadID := downloadCmd.String("id", "", "File ID")
		downloadCmd.Parse(args[1:])

		cmd.DownloadFile(*downloadID)

	case "convert":
		dir := convertCmd.String("dir", "", "Directory path")
		removeMd := convertCmd.Bool("removemd", false, "Remove markdown files after conversion")
		convertCmd.Parse(args[1:])

		if *dir == "" {
			fmt.Println("Error: directory path is required")
			os.Exit(1)
		}
		cmd.ConvertDirectory(*dir, *removeMd)

	default:
		fmt.Println("Error: command not found")
	}
}
