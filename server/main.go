package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sogladev/golang-terminal-downloader/pkg/manifest"
)

// ThrottledReader wraps an io.ReadSeeker and throttles the data being read
type ThrottledReader struct {
	reader   io.ReadSeeker
	interval time.Duration
	chunk    int
}

// Read reads data in chunks and introduces a delay between reads
func (t *ThrottledReader) Read(p []byte) (int, error) {
	if len(p) > t.chunk {
		p = p[:t.chunk] // Limit the read size to the defined chunk size
	}
	n, err := t.reader.Read(p)
	if n > 0 {
		time.Sleep(t.interval) // Simulate bandwidth delay
	}
	return n, err
}

// Seek sets the offset for the next Read operation
func (t *ThrottledReader) Seek(offset int64, whence int) (int64, error) {
	return t.reader.Seek(offset, whence)
}

func main() {
	// Add command-line flag for throttle interval
	interval := flag.Int("interval", 10, "ms delay per chunk")
	// Generate a manifest file for the input directory
	createManifest := flag.Bool("create-manifest", false, "Generate manifest.json before starting the server")
	filesDir := flag.String("files", "files", "Directory containing the files to process")
	baseURL := flag.String("url", "http://localhost:8080/", "Base URL for file download links")
	version := flag.String("version", "1.0", "Manifest version")

	flag.Parse()

	if *createManifest {
		fmt.Println("Generating manifest...")
		err := manifest.GenerateManifest(*filesDir, *baseURL, *version)
		if err != nil {
			log.Fatalf("Error generating manifest: %v", err)
		}
		fmt.Println("Manifest generated successfully.")
		return // Exit after generating the manifest
	}

	// Custom handler to throttle file downloads
	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Throttling request for: %s\n", r.URL.Path)

		// Open the requested file
		filePath := r.URL.Path[1:]
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		// Wrap the file in a ThrottledReader
		throttledReader := &ThrottledReader{
			reader:   file,
			interval: time.Duration(*interval) * time.Millisecond, // ms delay per chunk
			chunk:    1024,                                        // 1KB per chunk
		}

		// Serve the content using the throttled reader
		http.ServeContent(w, r, filePath, time.Now(), throttledReader)
	})

	// Fallback handler for all other requests
	http.Handle("/", http.FileServer(http.Dir("./")))

	// Start the server
	port := "8080"
	log.Printf("Starting local server on http://localhost:%s/", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
