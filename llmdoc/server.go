package llmdoc

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

type SlugReader interface {
	Read(slug string) (string, error)
}

type FileReader struct {
	// Directory where the markdown files are stored
	Dir string
}

func (fr FileReader) Read(slug string) (string, error) {
	slugPath := filepath.Join(fr.Dir, slug+".md")

	f, err := os.Open(slugPath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func IndexHandler(t *template.Template) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}, nil
}

func ArticleHandler(sr SlugReader, t *template.Template) (http.HandlerFunc, error) {
	mdRenderer := goldmark.New(
		goldmark.WithExtensions(),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		if strings.Contains(slug, ".md") {
			slug = strings.Split(slug, ".")[0]
			article, err := sr.Read(slug)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			_, err = w.Write([]byte(article))
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
			return
		}

		articleMarkdown, err := sr.Read(slug)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		err = mdRenderer.Convert([]byte(articleMarkdown), w)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}, nil
}
