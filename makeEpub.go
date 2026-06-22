package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"wp-to-epub/internal/scraper"

	"github.com/bmaupin/go-epub"
)

func makeEpub(volumes []scraper.Volume) error {
	for _, v := range volumes {
		title := v.VolumeTitle
		e := epub.NewEpub(title)
		cover, err := e.AddImage(v.VolumeCover, strings.ReplaceAll(title, " ", "-")+"-cover.png")
		if err != nil {
			return err
		}
		e.SetCover(cover, "")
		for _, c := range v.Chapters {
			chapterBody := strings.Join(c.ChapterBody, "\n")
			for _, i := range c.ChapterImages {
				image, err := e.AddImage(i, getImageName(i))
				if err != nil {
					return err
				}
				chapterBody = strings.Replace(chapterBody, i, image, 1)
			}
			chapterBody = fmt.Sprintf("<h1>%s</h1>\n", c.ChapterTitle) + chapterBody
			e.AddSection(chapterBody, c.ChapterTitle, "", "")
		}
		fp := filepath.Join("epub", title+".epub")
		err = e.Write(fp)
		if err != nil {
			return err
		}
	}
	return nil
}

func getImageName(path string) string {
	pathBits := strings.Split(path, "/")
	imgName := pathBits[len(pathBits)-1]
	imgName = strings.ReplaceAll(imgName, " ", "-")
	return imgName
}
