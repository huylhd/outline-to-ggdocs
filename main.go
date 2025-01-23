package main

import (
	"flag"
	"fmt"
	"os"
	"outline-to-ggdocs/cmd"
	"outline-to-ggdocs/config"
	"outline-to-ggdocs/utils"
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
		fmt.Println("  automate - Automate the entire process")
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
		utils.LogError("command is required")
		flag.Usage()
		os.Exit(1)
	}

	outlineApiKey := config.AppConfig.OutlineApiKey

	if outlineApiKey == "" {
		utils.LogError("OUTLINE_API_KEY is required")
		os.Exit(1)
	}

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
	automateCmd := flag.NewFlagSet("automate", flag.ExitOnError)
	fetchCmd := flag.NewFlagSet("fetchFileOperations", flag.ExitOnError)

	switch args[0] {
	case "list":
		page := listCmd.Int("page", 1, "Page number, starts from 1")
		listCmd.Parse(args[1:])

		cmd.ListCollectionsCommand(*page)

	case "export":
		exportID := exportCmd.String("id", "", "Collection ID")
		exportCmd.Parse(args[1:])

		cmd.ExportCollectionsCommand(*exportID)

	case "download":
		downloadID := downloadCmd.String("id", "", "File ID")
		downloadCmd.Parse(args[1:])

		cmd.DownloadFile(*downloadID, "./")

	case "convert":
		dir := convertCmd.String("dir", "", "Directory path")
		removeMd := convertCmd.Bool("remove-md", false, "Remove markdown files after conversion")
		convertCmd.Parse(args[1:])

		if *dir == "" {
			utils.LogError("directory path is required")
			os.Exit(1)
		}
		cmd.ConvertDirectory(*dir, *removeMd)

	case "automate":
		page := automateCmd.Int("page", 0, "Page number, starts from 1 (leave empty to fetch all collections)")
		fromStep := automateCmd.Int("from-step", 0, "Step to start from. 0 = list, 1 = export, 2 = download and convert")
		toStep := automateCmd.Int("to-step", 2, "Step to end at. 0 = list, 1 = export, 2 = download and convert")
		automateCmd.Parse(args[1:])

		fromAutomateStep := cmd.AutomateStep(*fromStep)
		toAutomateStep := cmd.AutomateStep(*toStep)
		if fromAutomateStep < cmd.STEP_LIST || fromAutomateStep > cmd.STEP_DOWNLOAD_UNZIP_CONVERT {
			utils.LogError("invalid from-step value")
			os.Exit(1)
		}
		if toAutomateStep < cmd.STEP_LIST || toAutomateStep > cmd.STEP_DOWNLOAD_UNZIP_CONVERT {
			utils.LogError("invalid to-step value")
			os.Exit(1)
		}

		cmd.AutomateCommand(*page, fromAutomateStep, toAutomateStep)

	case "fetchFileOperations":
		fetchCmd.Parse(args[1:])
		cmd.FetchFileOperations()

	default:
		utils.LogError("unknown command")
		flag.Usage()
		os.Exit(1)
	}
}
