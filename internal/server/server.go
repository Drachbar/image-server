package server

import (
	"io/fs"
	"log"
	"net/http"

	"image-server/internal/assets"
)

type Config struct {
	UploadDir  string
	Port       int
	BaseURL    string
	APIKeyFile string
}

type Server struct {
	config  Config
	apiKeys map[string]string
}

func New(cfg Config) (*Server, error) {
	s := &Server{
		config:  cfg,
		apiKeys: make(map[string]string),
	}
	if err := s.loadAPIKeys(cfg.APIKeyFile); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", chainMiddleware(s.uploadHandler, withCors, s.withAuth))
	mux.HandleFunc("/delete", chainMiddleware(s.deleteHandler, withCors, s.withAuth))
	mux.HandleFunc("/healthcheck", chainMiddleware(s.healthCheckHandler, withCors))
	mux.HandleFunc("/api/images", chainMiddleware(s.imagesAPIHandler, withCors))
	mux.HandleFunc("/api/apps", chainMiddleware(s.appsAPIHandler, withCors))
	mux.HandleFunc("/gallery", chainMiddleware(func(w http.ResponseWriter, r *http.Request) {
		content, _ := assets.FS.ReadFile("public/index.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	}, withCors))
	mux.HandleFunc("/gallery/", chainMiddleware(func(w http.ResponseWriter, r *http.Request) {
		content, _ := assets.FS.ReadFile("public/app.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	}, withCors))

	sub, err := fs.Sub(assets.FS, "public")
	if err != nil {
		log.Fatalf("Kunde inte läsa public-katalog: %v", err)
	}
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.FS(sub))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(s.config.UploadDir))))
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		content, _ := assets.FS.ReadFile("public/robots.txt")
		w.Header().Set("Content-Type", "text/plain")
		w.Write(content)
	})

	return mux
}
