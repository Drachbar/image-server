package server

import (
	"encoding/json"
	"net/http"
	"os"
)

func (s *Server) getAppFromAPIKey(r *http.Request) (string, bool) {
	key := r.Header.Get(apiKeyHeader)
	app, exists := s.apiKeys[key]
	return app, exists
}

func (s *Server) loadAPIKeys(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&s.apiKeys)
}
