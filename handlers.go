package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func collectImages(filterApp string) []ImageEntry {
	var all []ImageEntry

	appEntries, err := os.ReadDir(uploadDir)
	if err != nil {
		return all
	}
	sort.Slice(appEntries, func(i, j int) bool {
		return appEntries[i].Name() < appEntries[j].Name()
	})

	for _, appEntry := range appEntries {
		if !appEntry.IsDir() {
			continue
		}
		app := appEntry.Name()
		if filterApp != "" && app != filterApp {
			continue
		}

		dir1Entries, _ := os.ReadDir(filepath.Join(uploadDir, app))
		sort.Slice(dir1Entries, func(i, j int) bool { return dir1Entries[i].Name() < dir1Entries[j].Name() })
		for _, dir1 := range dir1Entries {
			if !dir1.IsDir() {
				continue
			}
			dir2Entries, _ := os.ReadDir(filepath.Join(uploadDir, app, dir1.Name()))
			sort.Slice(dir2Entries, func(i, j int) bool { return dir2Entries[i].Name() < dir2Entries[j].Name() })
			for _, dir2 := range dir2Entries {
				if !dir2.IsDir() {
					continue
				}
				files, _ := os.ReadDir(filepath.Join(uploadDir, app, dir1.Name(), dir2.Name()))
				sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
				for _, f := range files {
					if !f.IsDir() {
						all = append(all, ImageEntry{
							URL: fmt.Sprintf("%s/%s/%s/%s/%s", baseUrl, app, dir1.Name(), dir2.Name(), f.Name()),
							App: app,
						})
					}
				}
			}
		}
	}
	return all
}

func appsAPIHandler(w http.ResponseWriter, r *http.Request) {
	appEntries, err := os.ReadDir(uploadDir)
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
		images := collectImages(app)
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

func imagesAPIHandler(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	all := collectImages(app)
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
