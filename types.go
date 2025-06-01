package main

const apiKeyHeader = "X-API-Key"

type HealthResponse struct {
	Status string `json:"status"`
}

type UploadResponse struct {
	URL string `json:"url"`
}
