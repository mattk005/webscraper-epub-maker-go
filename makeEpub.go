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
		e.SetCover(cover, "")
		for _, c := range v.Chapters {
			chapterBody := strings.Join(c.ChapterBody, "\n")
			chapterBody = fmt.Sprintf("<h1>%s</h1>\n", title) + chapterBody
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
