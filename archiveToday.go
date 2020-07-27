package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ArchiveTodayMemento struct {
	Time  time.Time
	Url   string
	First bool
	Last  bool
}

type ArchiveTodayTimemap struct {
	Mementos []ArchiveTodayMemento
	Url      string
}

func ParseArchiveTodayResponse(response *http.Response) (mementos []ArchiveTodayMemento, err error) {
	scanner := bufio.NewScanner(response.Body)

	// iterate over lines, extracting mementos
	// line format:
	// <http://archive.md/20181206183642/https://cako.io/>; rel="first last memento"; datetime="Thu, 06 Dec 2018 18:36:42 GMT",
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), ",\n")

		// skip original url or timemap url
		if strings.Index(line, `rel="original"`) != -1 {
			continue
		}
		if strings.Index(line, `rel="self"`) != -1 {
			continue
		}

		// split key-value parts at semi-colon
		parts := strings.Split(line, "; ")

		m := ArchiveTodayMemento{}

		for _, p := range parts {
			if strings.Index(p, "http") == 0 {
				m.Url = strings.Trim(p, "<>")
			}
			if strings.Index(p, "rel") == 0 {
				if strings.Index(p, "first") != -1 {
					m.First = true
				}
				if strings.Index(p, "last") != -1 {
					m.Last = true
				}
			}
			if strings.Index(p, "datetime") == 0 {
				kv := strings.Split(p, `"`)
				t, err := time.Parse(time.RFC1123, kv[1])
				if err != nil {
					return mementos, err
				}
				m.Time = t
			}
		}

		mementos = append(mementos, m)
	}

	return
}

func GetArchiveTodayTimemap(uri string) (timemap ArchiveTodayTimemap, err error) {
	timemapUri := fmt.Sprintf("https://archive.today/timemap/%s", uri)

	r, err := httpClient.Get(timemapUri)
	if err != nil {
		return timemap, err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return timemap, errors.New(fmt.Sprintf("GetArchiveTodayTimemap error: %s", r.Status))
	}
	m, err := ParseArchiveTodayResponse(r)
	timemap.Mementos = m
	timemap.Url = uri

	return timemap, err
}
