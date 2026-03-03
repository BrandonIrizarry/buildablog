package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/environment"
)

//go:embed gohtml/*.gohtml
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var tpl *template.Template

type config struct {
	environment.Env
	handler http.HandlerFunc
}

func main() {
	// Set up logging.
	logFile, err := os.OpenFile("buildablog.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("can't open logfile: %v", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load environment variables. Embed these into the config
	// struct so that the handlers can use them.
	env, err := environment.New()
	if err != nil {
		log.Fatal(err)
	}

	// cfg is for passing state to the various handlers.
	cfg := config{
		Env: env,
	}

	// Parse the templates up front.
	funcMap := template.FuncMap{
		"dec": func(value int) int {
			return value - 1
		},
		"humanReadable": func(t time.Time) string {
			return t.Format(time.DateOnly)
		},
	}

	tpl, err = template.New("global").Funcs(funcMap).ParseFS(templateFS, "gohtml/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	// Set up the server.
	//
	// The various handlers are named after their endpoints in a
	// more or less obvious manner.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts/{date}", cfg.getPostsDate)
	mux.HandleFunc("GET /{$}", cfg.getIndex)
	mux.HandleFunc("GET /posts", cfg.getPosts)
	mux.HandleFunc("GET /tags", cfg.getTags)
	mux.HandleFunc("GET /projects", cfg.getProjects)
	mux.HandleFunc("GET /rss", cfg.getRSS)

	// Static assets (CSS files etc.)
	mux.Handle("/static/", http.FileServerFS(staticFS))

	// The server assumes the responsibility of serving any assets
	// defined locally inside the blog repo itself.
	blogAssetsDir := fmt.Sprintf("%s/assets", cfg.BlogDir)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(blogAssetsDir))))

	// Launch the server.
	log.Print("Killing any previous server instance; starting server on port 3030")

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}
