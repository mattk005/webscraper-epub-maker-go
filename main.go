package main

import (
	"fmt"
	"net/url"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		fmt.Println(err)
	}
	volumesURLs, err := getVolumeURLs(&cfg)
	if err != nil {
		fmt.Println(err)
	}
	volumes := make([]volume, 0)
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
		vol := volume{
			volumeURL:   _URL.String(),
			chapterURLs: chaptersURLs,
		}
		volumes = append(volumes, vol)
	}
	for i := range volumes {
		volumes[i].parseChapters()
		// volumes[i].getCover()
	}
}
