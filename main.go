package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/isacben/llmdoc/llmdoc"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "llmdoc\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "llmdoc [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	port := flag.Int("port", 8888, "Port for the webserver")
	flag.Parse()

	//mux := http.NewServeMux()

	articleReader := llmdoc.FileReader{
		Dir: "../posts",
	}

	indexTemplate, err := template.ParseFiles(
		"views/index.html",
		"views/header.html",
		"views/styles.html",
		"views/sidebar.html",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load index.html file: %v\n", err)
		os.Exit(1)
	}

	articleTemplate, err := template.ParseFiles(
		"views/article.html",
		"views/header.html",
		"views/styles.html",
		"views/sidebar.html",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load article.html file: %v\n", err)
		os.Exit(1)
	}

	page := struct {
		Sidebar string
		Header  string
	}{}

	sidebarMarkdown, err := articleReader.Read("sidebar")
	if err != nil {
		fmt.Println("error: could not read sidebar")
		return
	}

	page.Sidebar = sidebarMarkdown

	indexHandler, err := llmdoc.IndexHandler(page, indexTemplate)
	if err != nil {
		log.Fatal("error: failed to run server:", err)
	}

	articleHandler, err := llmdoc.ArticleHandler(articleReader, articleTemplate)
	if err != nil {
		log.Fatal("error: failed to run server:", err)
	}

	//fs := http.FileServer(http.Dir("../posts/static"))
	//mux.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("GET /{slug...}", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.Trim(r.PathValue("slug"), "/")

		if slug == "" {
			// Handle root path
			indexHandler(w, r)
			return
		}

		// Handle article paths
		articleHandler(w, r)
	})

	log.Printf("Deploying server. Navigate to http://127.0.0.1:%v\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}
