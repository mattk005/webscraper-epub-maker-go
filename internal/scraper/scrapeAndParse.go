package scraper

import (
	"fmt"
	"net/url"
)

func ScrapeAndParse() []Volume {
	cfg, err := getConfig()
	if err != nil {
		fmt.Println(err)
	}
	volumesURLs, err := getVolumeURLs(&cfg)
	if err != nil {
		fmt.Println(err)
	}
	volumes := make([]Volume, 0)
	for _, URL := range volumesURLs {
		// fmt.Println(URL)
		_URL, err := url.Parse(URL)
		if err != nil {
			fmt.Println(err)
		}
		chaptersURLs, err := getChapterURLs(_URL, &cfg)
		if err != nil {
			fmt.Println(err)
		}
		vol := Volume{
			VolumeURL:   _URL.String(),
			chapterURLs: chaptersURLs,
		}
		volumes = append(volumes, vol)
	}
	for i := range volumes {
		fmt.Printf("Parsing volumes - loop: %v\n", i+1)
		volumes[i].parse()
		// break // For testing! there's alot of time.After(500 * time.Millisecond) and like 400 chapters
	}
	return volumes
}
