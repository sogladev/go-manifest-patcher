package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	// Serve files in the "server" directory
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	// Run the server
	port := "8080"
	log.Printf("Starting local server on http://localhost:%s/", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Simulate Delays
	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Simulate delay
		http.ServeFile(w, r, r.URL.Path[1:])
	})
}
