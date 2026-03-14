package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

func appFromContext(r *http.Request) string {
	return r.Context().Value(appContextKey).(string)
}

func (s *Server) appsAPIHandler(w http.ResponseWriter, r *http.Request) {
	appEntries, err := os.ReadDir(s.config.UploadDir)
	if err != nil {
		http.Error(w, "Kunde inte läsa bildkatalog", http.StatusInternalServerError)
		return
	}
	sort.Slice(appEntries, func(i, j int) bool {
		return appEntries[i].Name() < appEntries[j].Name()
	})

	var apps []AppEntry
	for _, appEntry := range appEntries {
		if !appEntry.IsDir() {
			continue
		}
		app := appEntry.Name()
		images := s.collectImages(app)
		if len(images) == 0 {
			continue
		}
		apps = append(apps, AppEntry{
			Name:      app,
			Thumbnail: images[0].URL,
			Count:     len(images),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apps); err != nil {
		log.Printf("Fel vid apps-svar: %v", err)
	}
}

func (s *Server) imagesAPIHandler(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	all := s.collectImages(app)
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}

	var page []ImageEntry
	if offset < len(all) {
		page = all[offset:end]
	} else {
		page = []ImageEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ImagesResponse{Images: page, HasMore: end < len(all)}); err != nil {
		log.Printf("Fel vid galleri-svar: %v", err)
	}
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(HealthResponse{Status: "ok"}); err != nil {
		log.Printf("Fel vid healthcheck-svar: %v", err)
	}
}

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
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

	fileURL, err := s.saveFile(file, handler.Filename, app)
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

func (s *Server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	app := appFromContext(r)

	file := r.URL.Query().Get("filename")
	if file == "" {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}

	if err := s.deleteFile(file, app); err != nil {
		log.Printf("Fel vid borttagning: %v", err)
		http.Error(w, "Kunde inte ta bort filen", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Filen är borttagen"))
}
