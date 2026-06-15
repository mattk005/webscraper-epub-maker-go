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
		fmt.Printf("Parsing volumes - loop: %v\n", i)
		volumes[i].parse()
		if i == 1 {
			break
		}
		// volumes[i].getCover()
	}
	for _, volume := range volumes {
		fmt.Printf("Volume Name: %s \n", volume.volumeTitle)
		fmt.Printf("Volume Cover URL: %s \n", volume.volumeCover)
		fmt.Printf("Volume length (p tags ): %d \n", len(volume.chapterURLs))
	}
}
