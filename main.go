package main

import (
	"flag"
	"fmt"
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
}

func main() {
	updateExisting := flag.Bool("update", true, "Update already-archived pages")
	includePages := flag.Bool("includePages", true, "Include Pages in archive")

	archiveToday := flag.Bool("archiveToday", true, "Archive to archive.today")
	wayback := flag.Bool("wayback", true, "Archive to Wayback Machine/archive.org")

	flag.Parse()

	var uris []string

	posts, err := GetGhostPosts(GHOST_BASE_URI, GHOST_API_KEY)
	if err != nil {
		panic(err)
	}

	for _, u := range posts.Posts {
		uris = append(uris, u.Url)
	}

	if *includePages {
		pages, err := GetGhostPages(GHOST_BASE_URI, GHOST_API_KEY)
		if err != nil {
			panic(err)
		}
		for _, u := range pages.Pages {
			uris = append(uris, u.Url)
		}
	}

	if *wayback {
		RunWayback(uris, *updateExisting)
	}
	if *archiveToday {
		RunArchiveToday(uris, *updateExisting)
	}
}

// Blocking call to archive.today for uris; prints status in order
func RunArchiveToday(uris []string, updateExisting bool) {
	statusChs := make([]chan string, len(uris))
	for i, u := range uris {
		statusChs[i] = make(chan string)
		go ArchiveTodayRoutine(u, updateExisting, statusChs[i])
	}

	for _, u := range statusChs {
		fmt.Println(<-u)
	}
}

// Blocking call to Wayback for uris; prints status in order
func RunWayback(uris []string, updateExisting bool) {
	statusChs := make([]chan string, len(uris))

	for i, u := range uris {
		go WaybackRoutine(u, updateExisting, statusChs[i])
	}

	for _, u := range statusChs {
		fmt.Println(<-u)
	}
}
