package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type WaybackSnapshot struct {
	Available bool
	Timestamp int64 `json:",string"`
	Status    int64 `json:",string"`
	Url       string
}

func (r *WaybackSnapshot) Time() time.Time {
	return time.Unix(r.Timestamp, 0)
}

type WaybackAvailableResponse struct {
	Archived_snapshots struct {
		Closest *WaybackSnapshot
	}
	Url string
}

func (r *WaybackAvailableResponse) IsAvailable() bool {
	return r.Archived_snapshots.Closest != nil && r.Archived_snapshots.Closest.Available
}

func GetWaybackAvailable(uri string) (available WaybackAvailableResponse, err error) {
	waybackUri := fmt.Sprintf("https://archive.org/wayback/available?url=%s", uri)

	r, err := httpClient.Get(waybackUri)
	if err != nil {
		return available, err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return available, errors.New(fmt.Sprintf("GetWaybackAvailable error: %s", r.Status))
	}

	d := json.NewDecoder(r.Body)
	err = d.Decode(&available)
	return available, err
}

func SaveWayback(uri string) (response *http.Response, err error) {
	waybackUri := fmt.Sprintf("https://web.archive.org/save/%s", uri)
	return httpClient.Get(waybackUri)
}
