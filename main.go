package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sub2api/sub2api/handler"
)

const (
	defaultPort    = 8080
	defaultHost    = "0.0.0.0"
	appName        = "sub2api"
	appVersion     = "1.0.0"
)

func main() {
	// Command-line flags
	port := flag.Int("port", getEnvInt("PORT", defaultPort), "Port to listen on")
	host := flag.String("host", getEnv("HOST", defaultHost), "Host address to bind")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/sub", handler.SubHandler)
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/", handler.IndexHandler)

	log.Printf("Starting %s v%s on %s", appName, appVersion, addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an integer environment variable or returns a default value.
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: invalid integer value for %s, using default %d", key, defaultValue)
	}
	return defaultValue
}
