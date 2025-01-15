# Outline to Google Docs CLI Tool

This tool supports listing collections, exporting them, downloading files, and converting Markdown files to .docx

## Commands
- `list`: Retrieve all collections from Outline.
- `export`: Export specific collections or all collections (`markdown` format)
- `download`: Download exported file by ID.
- `convert`: Convert all Markdown files in a directory to .docx format

## Prerequisites
### 1. Go
If you need to build the project from source, ensure that Go is installed on your system.

### 2. Pandoc
To use the `convert` command, install [`pandoc`](https://pandoc.org/installing.html)
```bash
brew install pandoc
```

### 3. Environment variables
The OUTLINE_API_KEY is required to authenticate API requests to Outline.
```bash
export OUTLINE_API_KEY=api_key
```

## Usage
Use the existing binary in the repo
```bash
./outline-to-ggdocs <command> [flags]
```

Or build one
```bash
go build -o outline-to-ggdocs
```

Help command
```bash
outline-to-ggdocs --help
```