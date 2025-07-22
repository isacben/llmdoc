package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

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

	articleReader := llmdoc.FileReader{
		Dir: "../posts",
	}

	indexTemplate, err := template.ParseFiles("views/index.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load index.html file: %v\n", err)
		os.Exit(1)
	}

	articleTemplate, err := template.ParseFiles("views/article.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load article.html file: %v\n", err)
		os.Exit(1)
	}

	indexHandler, err := llmdoc.IndexHandler(indexTemplate)
	if err != nil {
		log.Fatal("error: failed to run server:", err)
	}

	articleHandler, err := llmdoc.ArticleHandler(articleReader, articleTemplate)
	if err != nil {
		log.Fatal("error: failed to run server:", err)
	}

	http.HandleFunc("GET /", indexHandler)
	http.HandleFunc("GET /{slug}", articleHandler)

	log.Printf("Deploying server. Navigate to http://127.0.0.1:%v\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}
