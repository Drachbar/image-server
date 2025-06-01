package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func saveFileWithApp(file io.Reader, filename, app string) (string, error) {
	hash := sha1.New()
	tee := io.TeeReader(file, hash)

	fileExt := filepath.Ext(filename)
	hashedBytes, err := io.ReadAll(tee)
	if err != nil {
		return "", err
	}

	hashSum := hash.Sum(nil)
	hashedSum := hex.EncodeToString(hashSum)

	dir1 := hashedSum[:2]
	dir2 := hashedSum[2:4]
	fullDir := filepath.Join(uploadDir, app, dir1, dir2)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", err
	}

	finalPath := filepath.Join(fullDir, hashedSum+fileExt)
	if err := os.WriteFile(finalPath, hashedBytes, 0644); err != nil {
		return "", err
	}

	// Publik URL
	return fmt.Sprintf("%s/%s/%s/%s/%s%s", baseUrl, app, dir1, dir2, hashedSum, fileExt), nil
}

func deleteFile(file, app string) error {
	ext := filepath.Ext(file)
	hash := strings.TrimSuffix(file, ext)

	if len(hash) < 4 {
		return fmt.Errorf("ogiltig hash")
	}

	dir1 := hash[:2]
	dir2 := hash[2:4]
	fullPath := filepath.Join(uploadDir, app, dir1, dir2, file)

	return os.Remove(fullPath)
}
