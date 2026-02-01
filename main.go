package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/xdg-go/mdash/internal/server"
)

func main() {
	var (
		port = flag.String("port", "", "Port to serve on (default: available port in 3000-3999)")
		dir  = flag.String("dir", ".", "Directory to serve markdown files from")
	)
	flag.Parse()

	// Check if directory exists
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist", *dir)
	}

	// Create listener
	var ln net.Listener
	var err error
	if *port != "" {
		ln, err = net.Listen("tcp", ":"+*port)
	} else {
		ln, err = listenOnAvailablePort(3000, 3999)
	}
	if err != nil {
		log.Fatalf("Could not listen: %v", err)
	}

	srv := server.New(*dir)

	httpServer := &http.Server{
		Handler:      srv,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Starting mdash serving %s on http://localhost:%d/\n", *dir, ln.Addr().(*net.TCPAddr).Port)
	log.Fatal(httpServer.Serve(ln))
}

// listenOnAvailablePort finds and binds to an available port in the given range.
// It starts at a random offset to reduce collisions between instances.
func listenOnAvailablePort(minPort, maxPort int) (net.Listener, error) {
	rangeSize := maxPort - minPort + 1
	start := rand.Intn(rangeSize)

	for i := range rangeSize {
		port := minPort + (start+i)%rangeSize
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			return ln, nil
		}
	}
	return nil, fmt.Errorf("no available port in range %d-%d", minPort, maxPort)
}
