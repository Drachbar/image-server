package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"image-server/internal/server"
)

func main() {
	var cfg server.Config
	flag.StringVar(&cfg.APIKeyFile, "apikeyfile", "/etc/apikeys.json", "Sökväg till API-nyckelfil")
	flag.StringVar(&cfg.UploadDir, "dir", "/var/www/images", "Sökväg för uppladdade bilder")
	flag.IntVar(&cfg.Port, "port", 8080, "Port att köra servern på")
	flag.StringVar(&cfg.BaseURL, "baseurl", "http://localhost:8080/images", "Publik URL-bas för bilder")
	flag.Parse()

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Kunde inte läsa API-nycklar: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Servern lyssnar på port %d...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(addr, srv.Routes()))
}
