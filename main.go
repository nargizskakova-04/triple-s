package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"triple-s/flag"
	"triple-s/handlers"
)

func main() {
	if err := flag.MyFlags(); err != nil {
		log.Fatalf("Flag error: %v\n", err)
	}

	port := *flag.Address
	baseDir := *flag.Dir

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.Mkdir(baseDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create base directory '%s' : %v\n", baseDir, err)
		}
	}

	mux := http.NewServeMux()

	mux.Handle("/", &handlers.BucketHandler{BaseDir: baseDir})
	mux.Handle("/{bucket}", &handlers.BucketHandler{BaseDir: baseDir})
	mux.Handle("/{bucket}/", &handlers.BucketHandler{BaseDir: baseDir})
	mux.Handle("/{bucket}/{object}", &handlers.ObjectHandler{BaseDir: baseDir})
	mux.Handle("/{bucket}/{object}/", &handlers.ObjectHandler{BaseDir: baseDir})

	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

