package main

import "net/http"

func checkAPIKey(r *http.Request) bool {
	return r.Header.Get(apiKeyHeader) == apiKey
}
