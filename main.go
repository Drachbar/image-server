package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	apiKey    string
	uploadDir string
	port      int
	baseUrl   string
)

const (
	apiKeyHeader = "X-API-Key"
)

func checkAPIKey(r *http.Request) bool {
	return r.Header.Get(apiKeyHeader) == apiKey
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-API-Key, Content-Type")

	if r.Method == http.MethodOptions {
		// Preflight request
		w.WriteHeader(http.StatusOK)
		return
	}

	if !checkAPIKey(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}

	defer file.Close()

	hash := sha1.New()
	tee := io.TeeReader(file, hash)

	hashedSum := ""
	tempFilePath := ""

	func() {
		fileExt := filepath.Ext(handler.Filename)
		hashedBytes, _ := io.ReadAll(tee)
		hashSum := hash.Sum(nil)
		hashedSum = hex.EncodeToString(hashSum)

		dir1 := hashedSum[:2]
		dir2 := hashedSum[2:4]

		fullDir := filepath.Join(uploadDir, dir1, dir2)
		err := os.MkdirAll(fullDir, 0755)

		if err != nil {
			http.Error(w, "Kunde inte spara fil", http.StatusInternalServerError)
			return
		}

		tempFilePath = filepath.Join(fullDir, hashedSum+fileExt)
		err = os.WriteFile(tempFilePath, hashedBytes, 0644)

		if err != nil {
			http.Error(w, "Kunde inte spara fil", http.StatusInternalServerError)
			return
		}
	}()

	fileURL := fmt.Sprintf("%s/%s/%s/%s%s", baseUrl, hashedSum[:2], hashedSum[2:4], hashedSum, filepath.Ext(handler.Filename))
	json.NewEncoder(w).Encode(map[string]string{"url": fileURL})
}

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

	http.HandleFunc("/upload", uploadHandler)
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Servern lyssnar på port %d...\n", port)
	http.ListenAndServe(addr, nil)
}
