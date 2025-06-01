package main

import "net/http"

func getAppFromAPIKey(r *http.Request) (string, bool) {
	key := r.Header.Get(apiKeyHeader)
	app, exists := apiKeys[key]
	return app, exists
}
