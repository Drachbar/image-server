package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(HealthResponse{Status: "ok"}); err != nil {
		log.Printf("Fel vid healthcheck-svar: %v", err)
	}
}

func appFromContext(r *http.Request) string {
	return r.Context().Value(appKey).(string)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	app := appFromContext(r)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileURL, err := saveFileWithApp(file, handler.Filename, app)
	if err != nil {
		log.Printf("Fel vid uppladdning: %v", err)
		http.Error(w, "Kunde inte spara fil", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(UploadResponse{URL: fileURL}); err != nil {
		log.Printf("Fel vid uppladdningssvar: %v", err)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	app := appFromContext(r)

	file := r.URL.Query().Get("filename")
	if file == "" {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}

	if err := deleteFile(file, app); err != nil {
		log.Printf("Fel vid borttagning: %v", err)
		http.Error(w, "Kunde inte ta bort filen", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Filen är borttagen"))
}
