package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/xdg-go/mdash/internal/server"
)

func main() {
	var (
		port = flag.String("port", "3000", "Port to serve on")
		dir  = flag.String("dir", ".", "Directory to serve markdown files from")
	)
	flag.Parse()

	// Check if directory exists
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist", *dir)
	}

	srv := server.New(*dir)

	httpServer := &http.Server{
		Addr:         ":" + *port,
		Handler:      srv,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Starting mdash serving %s on http://localhost:%s/\n", *dir, *port)
	log.Fatal(httpServer.ListenAndServe())
}
