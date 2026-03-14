package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var apiKeys map[string]string

func loadAPIKeys(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&apiKeys)
}

func init() {
	flag.StringVar(&apiKeyFile, "apikeyfile", "/etc/apikeys.json", "Sökväg till API-nyckelfil")
	flag.StringVar(&uploadDir, "dir", "/var/www/images", "Sökväg för uppladdade bilder")
	flag.IntVar(&port, "port", 8080, "Port att köra servern på")
	flag.StringVar(&baseUrl, "baseurl", "http://localhost:8080/images", "Publik URL-bas för bilder")
}

func main() {
	flag.Parse()

	if err := loadAPIKeys(apiKeyFile); err != nil {
		log.Fatalf("Kunde inte läsa API-nycklar: %v", err)
	}

	http.HandleFunc("/upload", chainMiddleware(uploadHandler, withCors, withAuth))
	http.HandleFunc("/delete", chainMiddleware(deleteHandler, withCors, withAuth))
	http.HandleFunc("/healthcheck", chainMiddleware(healthCheckHandler, withCors))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(uploadDir))))
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Servern lyssnar på port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
