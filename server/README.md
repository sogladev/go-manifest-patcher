# Test Server for Go Manifest Patcher

This is a test server application designed to simulate file downloads for the patcher download client. It provides manifest generation and throttled file serving capabilities.

## Features

- Manifest generation for files in a specified directory
- Throttled file downloads to simulate network conditions
- Simple HTTP server for hosting files
- Configurable download speeds

## Usage

### Directory Setup

1. Create a `files` directory (or specify custom with `-files` flag)
2. Add test files to the directory. You can create a 1MB test file using:
   ```bash
   dd if=/dev/zero of=files/testfile bs=1024 count=1024
   ```

### Running the Server

1. First, generate the manifest:
   ```bash
   go run main.go -create-manifest
   ```

2. Then start the server:
   ```bash
   go run main.go
   ```

> **Important**: Run the server from the server directory for correct file path resolution.

### Command Line Options

```bash
go run main.go --help

Usage:
  -create-manifest
        Generate manifest.json before starting the server
  -files string
        Directory containing the files to process (default "files")
  -interval int
        ms delay per chunk (default 10)
  -url string
        Base URL for file download links (default "http://localhost:8080/")
  -version string
        Manifest version (default "1.0")
```


### Workflow

1. Create test files in the `files` directory
2. Generate manifest using `-create-manifest` flag
3. Start the server
4. Use the downloader client to test against this server

The server will throttle downloads to simulate real-world conditions, useful for testing download progress indicators and resumption capabilities in the client.