# Outline to Google Docs CLI Tool

This tool supports listing collections, exporting them, downloading files, and converting Markdown files to .docx

## Commands
- `list`: Retrieve all collections from Outline.
- `export`: Export specific collections or all collections (`markdown` format)
- `download`: Download exported file by ID.
- `convert`: Convert all Markdown files in a directory to .docx format
- `automate`: Automate the entire process

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
<b>Option 1.</b> Use the prebuilt binary in the repo
```bash
./outline-to-ggdocs <command> [flags]
```

<b>Option 2.</b> Build the binary from source
```bash
go build -o outline-to-ggdocs
```

To see available commands, run:
```bash
outline-to-ggdocs --help
```
## The `automate` command
```bash
~ outline-to-ggdocs automate --help

Usage of automate:
  -from-step int
        Step to start from. 0 = list, 1 = exports, 2 = download and convert
  -page int
        Page number, starts from 1 (leave empty to fetch all collections)
```

### Steps
#### 1. List.
Fetches the collections.

- If the --page flag is specified, only the collections on that page will be fetched

- The result is saved into `collections.json`

#### 2. Exports
Read `collections.json` file and initiates export requests for collections.

- The operation data is saved into `fileOperations.json`

#### 3. Download and convert
Read `fileOperations.json`, check the statuses of export operations, and download the files into `/exports` directory.

- Unzip the files and converted all .md to .docx


<pre><b>💡 Tips:</b> 
Use "from-step" flag to skip a step. Useful when exporting operations takes a long time</pre>

#### Rate limit
Outline API has a rate limit of 50 calls / hour for the collections.export API. To avoid hitting the rate limit, it is recommended to use the `page` flag. Default is 50 limit each page