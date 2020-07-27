package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type GhostPost struct {
	Slug        string
	Id          string
	Uuid        string
	Title       string
	Visibility  string
	Created_at  time.Time
	Updated_at  time.Time
	Publised_at time.Time
}

type GhostPostsResponse struct {
	Posts []GhostPost
}

func GetGhostPosts(domain string, apiKey string) (posts GhostPostsResponse, err error) {
	uri := fmt.Sprintf("%s/ghost/api/v3/content/posts/?key=%s", domain, apiKey)

	r, err := httpClient.Get(uri)
	if err != nil {
		return posts, err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return posts, errors.New(fmt.Sprintf("GetGhostPosts error: %s", r.Status))
	}

	d := json.NewDecoder(r.Body)
	err = d.Decode(&posts)
	return posts, err
}

type GhostPagesResponse struct {
	Pages []GhostPost
}

func GetGhostPages(domain string, apiKey string) (pages GhostPagesResponse, err error) {
	uri := fmt.Sprintf("%s/ghost/api/v3/content/pages/?key=%s", domain, apiKey)

	r, err := httpClient.Get(uri)
	if err != nil {
		return pages, err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return pages, errors.New(fmt.Sprintf("GetGhostPages error: %s", r.Status))
	}

	d := json.NewDecoder(r.Body)
	err = d.Decode(&pages)
	return pages, err
}
