package server

const apiKeyHeader = "X-API-Key"

type HealthResponse struct {
	Status string `json:"status"`
}

type UploadResponse struct {
	URL string `json:"url"`
}

type ImageEntry struct {
	URL string `json:"url"`
	App string `json:"app"`
}

type ImagesResponse struct {
	Images  []ImageEntry `json:"images"`
	HasMore bool         `json:"hasMore"`
}

type AppEntry struct {
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Count     int    `json:"count"`
}
