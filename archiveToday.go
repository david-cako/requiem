package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const ARCHIVE_TODAY_WIP_POLLRATE = time.Second * 5

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

func (t *ArchiveTodayTimemap) Last() *ArchiveTodayMemento {
	for _, m := range t.Mementos {
		if m.Last {
			return &m
		}
	}
	return nil
}

func ParseArchiveTodayResponse(response *http.Response) (mementos []ArchiveTodayMemento, err error) {
	scanner := bufio.NewScanner(response.Body)

	// iterate over lines, extracting mementos
	// line format:
	// <http://archive.md/20181206183642/https://cako.io/>; rel="first last memento"; datetime="Thu, 06 Dec 2018 18:36:42 GMT",
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), ",\n")

		// skip original url, timegate or timemap url
		if strings.Index(line, `rel="original"`) != -1 {
			continue
		}
		if strings.Index(line, `rel="timegate"`) != -1 {
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
		// do not return error on 404, this means no entries
		if r.StatusCode == 404 {
			timemap.Url = uri
			return timemap, nil
		}
		return timemap, errors.New(r.Status)
	}
	m, err := ParseArchiveTodayResponse(r)
	timemap.Mementos = m
	timemap.Url = uri

	return timemap, err
}

func SaveArchiveToday(uri string) (archiveUri string, err error) {
	v := url.Values{}
	v.Add("url", uri)
	v.Add("always", "1")

	r, err := httpClient.PostForm("https://archive.today", v)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return "", errors.New(r.Status)
	}
	if r.Header.Get("refresh") == "" {
		return "", fmt.Errorf("No wip returned.")
	}

	for {
		// poll WIP url until 302
		wip, err := httpClient.Get(r.Header.Get("refresh"))
		if err != nil {
			return "", err
		}
		wip.Body.Close()

		if wip.StatusCode == 302 {
			// WIP Complete, Return location
			return wip.Header.Get("location"), nil
		} else if wip.StatusCode != 200 {
			return "", errors.New(wip.Status)
		}
		time.Sleep(ARCHIVE_TODAY_WIP_POLLRATE)
	}
}

func ArchiveTodayRoutine(uri string, updateExisting bool, statusCh chan<- string) {
	t, err := GetArchiveTodayTimemap(uri)
	if err != nil {
		statusCh <- fmt.Sprintf("ERROR GETTING ARCHIVE.TODAY TIMEMAP %s: %s", uri, err.Error())
	}

	if t.Last() == nil || updateExisting {
		_, err := SaveArchiveToday(uri)
		if err != nil {
			statusCh <- fmt.Sprintf("ARCHIVE.TODAY ERROR %s: %s", uri, err.Error())
		}
		statusCh <- fmt.Sprintf("archive.today success: %s", uri)
	} else {
		statusCh <- fmt.Sprintf("archive.today exists: %s", t.Last().Url)
	}
}
