package llmdoc

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type SlugReader interface {
	Read(slug string) (string, error)
}

type FileReader struct {
	// Directory where the markdown files are stored
	Dir string
}

type Page struct {
	Sidebar string
	Header  string
}

func (fr FileReader) Read(slug string) (string, error) {
	slugPath := filepath.Join(fr.Dir, slug+".md")

	//f, err := os.Open(slugPath)
	//if err != nil {
	//	return "", err
	//}

	//defer f.Close()

	//b, err := io.ReadAll(f)
	//if err != nil {
	//	return "", err
	//}

	// Check if file exists
	if _, err := os.Stat(slugPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", slugPath)
	}

	content, err := os.ReadFile(slugPath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return string(content), nil
}

func IndexHandler(page Page, t *template.Template) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := t.Execute(w, page)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}, nil
}

func ArticleHandler(sr SlugReader, t *template.Template) (http.HandlerFunc, error) {
	mdRenderer := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
		),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.Trim(r.PathValue("slug"), "/")
		if slug == "" {
			http.NotFound(w, r)
			return
		}

		if strings.Contains(slug, ".md") {
			slug = strings.Split(slug, ".")[0]
			article, err := sr.Read(slug)
			if err != nil {
				//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				http.NotFound(w, r)
				fmt.Println(err)
				return
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
			return
		}

		var buf bytes.Buffer
		err = mdRenderer.Convert([]byte(articleMarkdown), &buf)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		article := struct {
			Content template.HTML
		}{
			Content: template.HTML(buf.String()),
		}

		err = t.Execute(w, article)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}, nil
}
