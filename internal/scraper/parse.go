package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Volume struct {
	VolumeURL   string
	VolumeTitle string
	VolumeCover string
	chapterURLs []string
	Chapters    []Chapter
}

type Chapter struct {
	ChapterTitle string
	ChapterBody  []string
	bookmarkData BookmarkData
}

type BookmarkData struct {
	Thumbnail  string `json:"thumbnail"`
	Cover      string `json:"cover"`
	Title      string `json:"title"`
	StoryTitle string `json:"storyTitle"`
}

func (v *Volume) parse() {
	for i, chapterURL := range v.chapterURLs {
		chapterHTML, err := getHTML(chapterURL)
		if err != nil {
			fmt.Println(err)
		}
		reader := strings.NewReader(chapterHTML)
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		chap := Chapter{}
		chap.getTitle(doc)
		chap.getBody(doc)
		chap.getBookmarkData(doc)
		v.Chapters = append(v.Chapters, chap)
		if i == 0 {
			v.VolumeCover = v.Chapters[0].bookmarkData.Cover
			v.VolumeTitle = v.Chapters[0].bookmarkData.StoryTitle
			v.saveCover()
		}
		time.After(500 * time.Millisecond)
	}
	if len(v.Chapters) == 0 {
		fmt.Println("no chapters found (Cover/Title)")
	}
}

func (v *Volume) saveCover() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", v.Chapters[0].bookmarkData.Cover, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	filePath := filepath.Join("internal", "scraper", "covers", v.VolumeTitle+".png")
	filePath = strings.ReplaceAll(filePath, " ", "-")
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	v.VolumeCover = filePath
	return err
}

func (c *Chapter) getTitle(doc *goquery.Document) (string, error) {
	title := doc.Find("h1.chapter__title").First().Text()
	c.ChapterTitle = title
	return title, nil
}

func (c *Chapter) getBody(doc *goquery.Document) ([]string, error) {
	var results []string
	content := doc.Find("#chapter-content .chapter-formatting")
	firstP := content.Find("p").First()
	rest := firstP.NextUntil("hr")
	paragraphs := firstP.Union(rest)

	paragraphs.EachWithBreak(func(i int, s *goquery.Selection) bool {
		// Check if the current paragraph contains the specific "Previous" link
		// You can check for the text or the href attribute
		link := s.Find("a")
		if link.AttrOr("data-type", "") == "fcn_chapter" {
			return false // Break the EachWithBreak loop
		}
		if s.Is("hr") {
			return false
		}

		// Your existing cleanup logic
		s.RemoveAttr("id")
		s.RemoveAttr("data-paragraph-id")

		s.Find("span").Each(func(j int, span *goquery.Selection) {
			spanHTML, err := span.Html()
			if err != nil {
				fmt.Println(spanHTML)
			}
			span.ReplaceWithHtml(spanHTML)
		})

		outer, _ := goquery.OuterHtml(s)
		results = append(results, outer)

		return true // Continue to next
	})

	// paragraphs.Each(func(i int, s *goquery.Selection) {
	// 	s.RemoveAttr("id")
	// 	s.RemoveAttr("data-paragraph-id")
	//
	// 	outer, err := goquery.OuterHtml(s)
	// 	if err == nil {
	// 		results = append(results, outer)
	// 	}
	// })
	c.ChapterBody = results
	return results, nil
}

func (c *Chapter) getBookmarkData(doc *goquery.Document) (BookmarkData, error) {
	var data BookmarkData
	selection := doc.Find("#fictioneer-bookmark-data")
	jsonString := selection.Text()
	err := json.Unmarshal([]byte(jsonString), &data)
	c.bookmarkData = data
	return data, err
}
