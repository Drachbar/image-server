package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func init() {
	flag.StringVar(&apiKey, "apikey", "", "API-nyckel som krävs för uppladdning")
	flag.StringVar(&uploadDir, "dir", "/var/www/images", "Sökväg för uppladdade bilder")
	flag.IntVar(&port, "port", 8080, "Port att köra servern på")
	flag.StringVar(&baseUrl, "baseurl", "http://localhost:8080/images", "Publik URL-bas för bilder")
}

func main() {
	flag.Parse()

	if apiKey == "" {
		log.Fatal("API-nyckel krävs, starta med -apikey")
	}

	http.HandleFunc("/upload", chainMiddleware(uploadHandler, withCors, withAuth))
	http.HandleFunc("/delete", chainMiddleware(deleteHandler, withCors, withAuth))
	http.HandleFunc("/healthcheck", chainMiddleware(healthCheckHandler, withCors))
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Servern lyssnar på port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
