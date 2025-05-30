package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func saveFile(file io.Reader, filename string) (string, error) {
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
	fullDir := filepath.Join(uploadDir, dir1, dir2)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", err
	}

	finalPath := filepath.Join(fullDir, hashedSum+fileExt)
	fmt.Println("Sparar till:", finalPath)

	if err := os.WriteFile(finalPath, hashedBytes, 0644); err != nil {
		return "", err
	}

	// Returnera publik URL
	return fmt.Sprintf("%s/%s/%s/%s%s", baseUrl, dir1, dir2, hashedSum, fileExt), nil
}
