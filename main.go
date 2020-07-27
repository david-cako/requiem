package main

import (
	"flag"
	"net/http"
	"os"
	"time"
)

var httpClient *http.Client

var GHOST_BASE_URI string
var GHOST_API_KEY string

func init() {
	GHOST_BASE_URI = os.Getenv("GHOST_BASE_URI")
	GHOST_API_KEY = os.Getenv("GHOST_API_KEY")
	if GHOST_BASE_URI == "" || GHOST_API_KEY == "" {
		panic("GHOST_BASE_URI or GHOST_API_KEY environment variable undefined.")
	}

	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	update := flag.Bool("update", true, "Update already-archived pages")
	archiveToday := flag.Bool("archiveToday", true, "Archive to archive.today")
	wayback := flag.Bool("wayback", true, "Archive to Wayback Machine/archive.org")

	flag.Parse()
}
