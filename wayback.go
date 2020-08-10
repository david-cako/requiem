package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type WaybackSnapshot struct {
	Available bool
	Timestamp string
	Status    int64 `json:",string"`
	Url       string
}

func (r *WaybackSnapshot) Time() (t time.Time, err error) {
	y, err := strconv.Atoi(r.Timestamp[0:4])
	if err != nil {
		return t, err
	}
	m, err := strconv.Atoi(r.Timestamp[4:6])
	if err != nil {
		return t, err
	}
	d, err := strconv.Atoi(r.Timestamp[6:8])
	if err != nil {
		return t, err
	}
	h, err := strconv.Atoi(r.Timestamp[8:10])
	if err != nil {
		return t, err
	}
	min, err := strconv.Atoi(r.Timestamp[10:12])
	if err != nil {
		return t, err
	}
	s, err := strconv.Atoi(r.Timestamp[12:14])
	if err != nil {
		return t, err
	}
	return time.Date(y, time.Month(m), d, h, min, s, 0, time.UTC), nil
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
		return available, errors.New(r.Status)
	}

	d := json.NewDecoder(r.Body)
	err = d.Decode(&available)
	return available, err
}

func SaveWayback(uri string) (archiveUri string, err error) {
	v := url.Values{}
	v.Add("url", uri)
	v.Add("capture_all", "on")
	r, err := httpClient.PostForm("https://web.archive.org/save/", v)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return "", errors.New(r.Status)
	}

	return fmt.Sprintf("https://web.archive.org%s", r.Header.Get("content-location")), nil
}

func WaybackRoutine(uri string, updateExisting bool, statusCh chan<- string) {
	a, err := GetWaybackAvailable(uri)
	if err != nil {
		statusCh <- fmt.Sprintf("ERROR CHECKING WAYBACK AVAILABLE %s: %s", uri, err.Error())
	}

	if !a.IsAvailable() || updateExisting {
		_, err := SaveWayback(uri)
		if err != nil {
			statusCh <- fmt.Sprintf("WAYBACK ERROR %s: %s", uri, err.Error())
		}
		statusCh <- fmt.Sprintf("Wayback success: %s", uri)
	} else {
		statusCh <- fmt.Sprintf("Wayback exists: %s", a.Url)
	}
}
