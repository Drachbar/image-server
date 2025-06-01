package main

import "net/http"

func checkAPIKey(r *http.Request) bool {
	header := r.Header.Get(apiKeyHeader)

	return header == apiKey
}
