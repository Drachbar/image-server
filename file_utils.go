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

func hashDirPath(app, hash string) string {
	return filepath.Join(uploadDir, app, hash[:2], hash[2:4])
}

func saveFileWithApp(file io.Reader, filename, app string) (string, error) {
	hash := sha1.New()
	tee := io.TeeReader(file, hash)

	fileExt := filepath.Ext(filename)
	hashedBytes, err := io.ReadAll(tee)
	if err != nil {
		return "", err
	}

	hashedSum := hex.EncodeToString(hash.Sum(nil))
	fullDir := hashDirPath(app, hashedSum)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", err
	}

	finalPath := filepath.Join(fullDir, hashedSum+fileExt)
	if err := os.WriteFile(finalPath, hashedBytes, 0644); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s/%s/%s%s", baseUrl, app, hashedSum[:2], hashedSum[2:4], hashedSum, fileExt), nil
}

func deleteFile(file, app string) error {
	ext := filepath.Ext(file)
	hash := strings.TrimSuffix(file, ext)

	if len(hash) < 4 {
		return fmt.Errorf("ogiltig hash")
	}

	fullPath := filepath.Join(hashDirPath(app, hash), file)
	return os.Remove(fullPath)
}
