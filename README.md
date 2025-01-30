# Go Manifest Patcher

A lightweight Go-based terminal patcher that uses a manifest to manage file updates. It displays a transaction overview, provides detailed progress, and only overwrites files listed in the manifest. It does not remove extra files. Designed for easy extension with minimal dependencies.

![downloader](images/downloader.gif)

## Features

- Manifest-based file synchronization
- Transaction overview before downloading
- Automatic updates from latest github release
- Cross-platform, windows and linux
- Smart file categorization:
  - Up-to-date files (skipped)
  - Outdated files (updated)
  - Missing files (downloaded)
  - Extra files detection (user-defined filter)
- Progress visualization with speed and ETA
- Support for both local and remote manifests

## Usage

### Basic Usage

```bash
go run main.go
```

By default, looks for `manifest.json` in the current directory.

### Command Line Options

```bash
go run main.go --help

Usage:
  -log-level string
        Set the log level (debug, info, warning, error) (default "info")
  -manifest string
        Path to manifest.json file or URL (e.g., http://localhost:8080/manifest.json) (default "manifest.json")
  -save-filter
        Save the default filter to filter.json and exit
  -skip-update
        Skip update check (useful for development)

----------------
```

### Transaction Overview

The downloader provides a detailed overview before executing downloads:

1. **Up-to-date files**: Files that match the manifest
2. **Outdated files**: Existing files that need updating
3. **Missing files**: New files to download
4. **Extra files**: Files not in manifest that are not ignored by custom filter
5. **Transaction summary**: Shows total download size and disk space impact

You'll be prompted to confirm before proceeding with downloads.

```bash
 go run main.go -manifest http://localhost:8080/manifest.json
    ____
   / __ )  ____ _   ____    ____   ___    _____
  / __  | / __ `/  / __ \  / __ \ / _ \  / ___/
 / /_/ / / /_/ /  / / / / / / / //  __/ / /
/_____/  \__,_/  /_/ /_/ /_/ /_/ \___/ /_/

Downloading manifest from: http://localhost:8080/manifest.json

Manifest Overview:
 Version: 1.0
 Up-to-date files:
  files/A.bin (Size: 1.0 MB)
  files/B.bin (Size: 1.0 MB)
  files/more/E.bin (Size: 1.0 MB)

 Outdated files (will be updated):
  files/C.bin (Current Size: 1.0 MB, New Size: 2.1 MB)

 Missing files (will be downloaded):
  files/D.bin (New Size: 2.1 MB)
  files/more/F.bin (New Size: 1.0 MB)

 Extra files (not in manifest):
  files/Z.bin (Size: 1.0 MB)

Transaction Summary:
 Installing/Updating: 3 files

Total size of inbound files is 5.2 MB. Need to download 5.2 MB.
After this operation, 4.1 MB of additional disk space will be used.
Is this ok [y/N]: y
[1/3] D.bin                     [--------------------] 100% (complete) 2.1 MB
[2/3] F.bin                     [--------------------] 100% (complete) 1.0 MB
[3/3] C.bin                     [--------------------] 100% (complete) 2.1 MB

--------------------------------------------------------------------------------
All files are up to date or successfully downloaded.
```

## Development

### Testing with Local Server

1. Clone both the downloader and server repositories
2. Start the test server (see [server README](./server/README.md))
   ```bash
   cd ../server
   go run main.go
   go run main.go
   ```

3. Run the downloader pointing to the test server's manifest:
   ```bash
   cd ../downloader
   go run main.go -manifest "http://localhost:8080/manifest.json"
   ```

### Manifest Format

The manifest.json should follow this structure:
```json
{
  "Version": "1.0",
  "Files": [
    {
      "Path": "path/to/file",
      "Hash": "file-hash",
      "Size": fileSize,
      "Custom": true,
      "URL": "url-to-file"
    },
}

```

## Testing

For development testing, use the companion test server application which provides:
- Local CDN simulation
- Manifest generation
- Throttled downloads
- Configurable file serving

See the [server README](./server/README.md) for more details on setting up a test environment.


## Creating the Demo GIF

This demo GIF was created by:
1. Recording a terminal session with [asciinema](https://asciinema.org).
2. Converting the `.cast` file to GIF using [agg](https://docs.asciinema.org/manual/agg/).
3. Optionally trimming or optimizing the resulting GIF with [ezgif.com](https://ezgif.com).

## License

See [LICENSE](LICENSE)
