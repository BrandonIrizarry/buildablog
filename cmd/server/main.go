package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// tpls maps template basenames to actual templates. This is so that
// we can parse all our templates up front, as opposed to parsing them
// on each request.
var tpls = make(map[string]*template.Template)

type rssConfig struct {
	flagLocal bool
	handler   http.HandlerFunc
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

	// Flags: if -flagLocal is specified, set RSS siteURL to
	// localhost:PORT.
	var flagLocal bool
	flag.BoolVar(&flagLocal, "local", false, "Whether we're serving from localhost")
	flag.Parse()

	// Parse the templates up front.
	funcMap := template.FuncMap{
		"dec": func(value int) int {
			return value - 1
		},
		"humanReadable": func(t time.Time) string {
			return t.Format(time.DateOnly)
		},
	}

	gohtmlFiles, err := filepath.Glob("gohtml/*")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range gohtmlFiles {
		log.Printf("Parsing template file '%s'", file)

		name := filepath.Base(file)
		tpl, err := template.New(name).Funcs(funcMap).ParseFiles(file, "html/nav.html")
		if err != nil {
			log.Fatal(err)
		}

		tpls[name] = tpl
	}

	// Set up the server.
	//
	// The various handlers are named after their endpoints in a
	// more or less obvious manner.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts/{date}", getPostsDate)
	mux.HandleFunc("GET /{$}", getIndex)
	mux.HandleFunc("GET /posts", getPosts)
	mux.HandleFunc("GET /tags", getTags)
	mux.HandleFunc("GET /projects", getProjects)

	// Serve the RSS feed. For now we need this "config" trick if
	// we're going to pass state to the handler now. Hopefully in
	// production we won't need anything like this.
	rssCfg := rssConfig{
		flagLocal: flagLocal,
	}

	mux.HandleFunc("GET /rss", rssCfg.getRSS)

	// Static assets (CSS files etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Launch the server.
	log.Print("Killing any previous server instance; starting server on port 3030")

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}
