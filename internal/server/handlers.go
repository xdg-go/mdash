package server

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PageData struct {
	Title   string
	Content template.HTML
	Path    string
}

type DirectoryData struct {
	Title   string
	Path    string
	Entries []DirEntry
}

type DirEntry struct {
	Name  string
	Path  string
	IsDir bool
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	if urlPath == "/" {
		urlPath = ""
	}

	fsPath := filepath.Join(s.baseDir, urlPath)

	// Security check: ensure path is within baseDir
	cleanPath, err := filepath.Abs(fsPath)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	baseAbs, err := filepath.Abs(s.baseDir)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if !strings.HasPrefix(cleanPath, baseAbs) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	info, err := os.Stat(cleanPath)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		s.serveDirectory(w, r, cleanPath, urlPath)
	} else if strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
		s.serveMarkdown(w, r, cleanPath, urlPath)
	} else {
		http.ServeFile(w, r, cleanPath)
	}
}

func (s *Server) serveDirectory(w http.ResponseWriter, r *http.Request, fsPath, urlPath string) {
	// Check if index.md exists in this directory
	indexPath := filepath.Join(fsPath, "index.md")
	if info, err := os.Stat(indexPath); err == nil && !info.IsDir() {
		// Serve index.md instead of directory listing
		s.serveMarkdown(w, r, indexPath, urlPath)
		return
	}

	entries, err := os.ReadDir(fsPath)
	if err != nil {
		http.Error(w, "Cannot read directory", http.StatusInternalServerError)
		return
	}

	var dirEntries []DirEntry

	// Add parent directory link if not at root
	if urlPath != "" {
		parent := filepath.Dir(urlPath)
		parentPath := "/"
		if parent != "." {
			parentPath = parent
		}
		dirEntries = append(dirEntries, DirEntry{
			Name:  "..",
			Path:  parentPath,
			IsDir: true,
		})
	}

	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") && !entry.IsDir() {
			continue
		}

		entryPath := "/"
		if urlPath != "" {
			entryPath = urlPath + "/"
		}
		entryPath += entry.Name()

		dirEntries = append(dirEntries, DirEntry{
			Name:  entry.Name(),
			Path:  entryPath,
			IsDir: entry.IsDir(),
		})
	}

	// Sort: directories first, then files
	sort.Slice(dirEntries, func(i, j int) bool {
		if dirEntries[i].Name == ".." {
			return true
		}
		if dirEntries[j].Name == ".." {
			return false
		}
		if dirEntries[i].IsDir != dirEntries[j].IsDir {
			return dirEntries[i].IsDir
		}
		return dirEntries[i].Name < dirEntries[j].Name
	})

	title := "Index"
	if urlPath != "" {
		title = urlPath
	}

	data := DirectoryData{
		Title:   title,
		Path:    urlPath,
		Entries: dirEntries,
	}

	tmpl, err := template.ParseFS(staticFS, "templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func (s *Server) serveMarkdown(w http.ResponseWriter, r *http.Request, fsPath, urlPath string) {
	content, err := os.ReadFile(fsPath)
	if err != nil {
		http.Error(w, "Cannot read file", http.StatusInternalServerError)
		return
	}

	html, err := s.renderer.Render(content)
	if err != nil {
		http.Error(w, "Markdown rendering error", http.StatusInternalServerError)
		return
	}

	title := filepath.Base(urlPath)
	if strings.HasSuffix(title, ".md") {
		title = title[:len(title)-3]
	}

	data := PageData{
		Title:   title,
		Content: template.HTML(html),
		Path:    urlPath,
	}

	tmpl, err := template.ParseFS(staticFS, "templates/document.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
