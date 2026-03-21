package server

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (s *Server) hashDirPath(app, hash string) string {
	return filepath.Join(s.config.UploadDir, app, hash[:2], hash[2:4])
}

func (s *Server) saveFile(file io.Reader, filename, app string) (string, error) {
	if err := os.MkdirAll(s.config.UploadDir, 0755); err != nil {
		return "", err
	}
	tmpFile, err := os.CreateTemp(s.config.UploadDir, "imgupload-*")
	if err != nil {
		return "", err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath) // no-op efter lyckad rename

	hash := sha1.New()
	if _, err := io.Copy(tmpFile, io.TeeReader(file, hash)); err != nil {
		tmpFile.Close()
		return "", err
	}
	if err := tmpFile.Close(); err != nil {
		return "", err
	}

	hashedSum := hex.EncodeToString(hash.Sum(nil))
	fileExt := filepath.Ext(filename)
	fullDir := s.hashDirPath(app, hashedSum)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", err
	}

	finalPath := filepath.Join(fullDir, hashedSum+fileExt)
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s/%s/%s%s", s.config.BaseURL, app, hashedSum[:2], hashedSum[2:4], hashedSum, fileExt), nil
}

func (s *Server) deleteFile(file, app string) error {
	ext := filepath.Ext(file)
	hash := strings.TrimSuffix(file, ext)

	if len(hash) < 4 {
		return fmt.Errorf("ogiltig hash")
	}

	fullPath := filepath.Join(s.hashDirPath(app, hash), file)
	return os.Remove(fullPath)
}

func (s *Server) collectImages(filterApp string) []ImageEntry {
	var all []ImageEntry

	appEntries, err := os.ReadDir(s.config.UploadDir)
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

		dir1Entries, _ := os.ReadDir(filepath.Join(s.config.UploadDir, app))
		sort.Slice(dir1Entries, func(i, j int) bool { return dir1Entries[i].Name() < dir1Entries[j].Name() })
		for _, dir1 := range dir1Entries {
			if !dir1.IsDir() {
				continue
			}
			dir2Entries, _ := os.ReadDir(filepath.Join(s.config.UploadDir, app, dir1.Name()))
			sort.Slice(dir2Entries, func(i, j int) bool { return dir2Entries[i].Name() < dir2Entries[j].Name() })
			for _, dir2 := range dir2Entries {
				if !dir2.IsDir() {
					continue
				}
				files, _ := os.ReadDir(filepath.Join(s.config.UploadDir, app, dir1.Name(), dir2.Name()))
				sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
				for _, f := range files {
					if !f.IsDir() {
						all = append(all, ImageEntry{
							URL: fmt.Sprintf("%s/%s/%s/%s/%s", s.config.BaseURL, app, dir1.Name(), dir2.Name(), f.Name()),
							App: app,
						})
					}
				}
			}
		}
	}
	return all
}
