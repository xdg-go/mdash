package server

import (
	"embed"
	"net/http"

	"github.com/xdg-go/mdash/internal/markdown"
)

//go:embed static templates
var staticFS embed.FS

// Server handles HTTP requests for markdown files
type Server struct {
	baseDir  string
	renderer *markdown.Renderer
	mux      *http.ServeMux
}

// New creates a new server instance
func New(baseDir string) *Server {
	s := &Server{
		baseDir:  baseDir,
		renderer: markdown.New(),
		mux:      http.NewServeMux(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Serve static files
	s.mux.Handle("/static/", http.FileServer(http.FS(staticFS)))

	// Handle all other requests
	s.mux.HandleFunc("/", s.handleRequest)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
