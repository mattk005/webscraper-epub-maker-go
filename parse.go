package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type volume struct {
	volumeURL   string
	volumeTitle string
	volumeCover string
	chapterURLs []string
	chapters    []chapter
}

type chapter struct {
	chapterTitle string
	chapterBody  []string
	bookmarkData BookmarkData
}

type BookmarkData struct {
	Thumbnail  string `json:"thumbnail"`
	Cover      string `json:"cover"`
	Title      string `json:"title"`
	StoryTitle string `json:"storyTitle"`
}

func (v *volume) parseChapters() {
	for _, chapterURL := range v.chapterURLs {
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
		chap := chapter{}
		chap.getTitle(doc)
		chap.getBody(doc)
		chap.getBookmarkData(doc)

		v.chapters = append(v.chapters, chap)
	}
	if len(v.chapters) > 0 {
		v.volumeCover = v.chapters[0].bookmarkData.Cover
	} else {
		fmt.Println("no chapters found (no cover data)")
	}
}

func (c *chapter) getTitle(doc *goquery.Document) (string, error) {
	title := doc.Find("h1.chapter__title").First().Text()
	c.chapterTitle = title
	return title, nil
}

func (c *chapter) getBody(doc *goquery.Document) ([]string, error) {
	var results []string
	content := doc.Find("#chapter-content .chapter-formatting")
	firstP := content.Find("p").First()
	rest := firstP.NextUntil("hr")
	paragraphs := firstP.Union(rest)

	paragraphs.Each(func(i int, s *goquery.Selection) {
		s.RemoveAttr("id")
		s.RemoveAttr("data-paragraph-id")

		s.Find("span").Each(func(j int, span *goquery.Selection) {
			spanHTML, err := span.Html()
			if err != nil {
				fmt.Println(spanHTML)
			}
			span.ReplaceWithHtml(spanHTML)
		})
		outer, err := goquery.OuterHtml(s)
		if err == nil {
			results = append(results, outer)
		}
	})
	c.chapterBody = results
	return results, nil
}

func (c *chapter) getBookmarkData(doc *goquery.Document) (BookmarkData, error) {
	var data BookmarkData
	selection := doc.Find("#fictioneer-bookmark-data")
	jsonString := selection.Text()
	err := json.Unmarshal([]byte(jsonString), &data)
	c.bookmarkData = data
	return data, err
}
