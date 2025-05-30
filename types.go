package main

var (
	apiKey    string
	uploadDir string
	port      int
	baseUrl   string
)

const apiKeyHeader = "X-API-Key"

type HealthResponse struct {
	Status string `json:"status"`
}

type UploadResponse struct {
	URL string `json:"url"`
}
