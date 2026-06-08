package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

func getVolumeURLs(cfg *config) ([]string, error) {
	baseURL := cfg.baseURL
	seriesURL := cfg.baseURL.ResolveReference(cfg.seriesPath)
	allVolumeURLs := make([]string, 0)
	volume1Found := false

	q := seriesURL.Query()
	for i := 1; ; i++ {
		q.Set("pg", fmt.Sprintf("%d", i))
		seriesURL.RawQuery = q.Encode()

		seriesPageHTML, err := getHTML(seriesURL.String())
		if err != nil {
			return []string{}, nil
		}
		volumeURLs, err := getURLsFromHTML(seriesPageHTML, baseURL)
		if err != nil {
			return []string{}, nil
		}

		for _, url := range volumeURLs {
			if strings.Contains(url, cfg.volumePattern) {
				// fmt.Println(url)
				allVolumeURLs = append(allVolumeURLs, url)
				if strings.Contains(url, cfg.volumePatternStop) {
					volume1Found = true
				}
			}
		}
		time.After(500 * time.Millisecond)
		if volume1Found || i > 4 {
			break
		}
	}
	return allVolumeURLs, nil
}

func getChapterURLs(volumeURL *url.URL, cfg *config) ([]string, error) {
	baseURL := cfg.baseURL
	allChapterURLs := make([]string, 0)

	volumePage, err := getHTML(volumeURL.String())
	if err != nil {
		return []string{}, nil
	}
	chapterURLs, err := getURLsFromHTML(volumePage, baseURL)
	if err != nil {
		return []string{}, nil
	}

	for _, url := range chapterURLs {
		if strings.Contains(url, cfg.chapterPattern) {
			// fmt.Println(url)
			allChapterURLs = append(allChapterURLs, url)
		}
	}
	return allChapterURLs, nil
}
