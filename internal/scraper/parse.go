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
	ChapterTitle  string
	ChapterBody   []string
	ChapterImages []string
	bookmarkData  BookmarkData
}

type BookmarkData struct {
	Thumbnail  string `json:"thumbnail"`
	Cover      string `json:"cover"`
	Title      string `json:"title"`
	StoryTitle string `json:"storyTitle"`
}

func (v *Volume) parse() error {
	for i, chapterURL := range v.chapterURLs {
		chapterHTML, err := getHTML(chapterURL)
		if err != nil {
			return err
		}
		reader := strings.NewReader(chapterHTML)
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			fmt.Println(err)
			return err
		}

		chap := Chapter{}

		_, err = chap.getTitle(doc)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = chap.getBody(doc)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = chap.getBookmarkData(doc)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// err = chap.saveIllustration()
		// if err != nil {
		// 	fmt.Println(err)
		// 	return err
		// }
		v.Chapters = append(v.Chapters, chap)
		if i == 0 {
			v.VolumeCover = v.Chapters[0].bookmarkData.Cover
			v.VolumeTitle = v.Chapters[0].bookmarkData.StoryTitle
			err = v.saveCover()
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		time.After(500 * time.Millisecond)
	}
	if len(v.Chapters) == 0 {
		fmt.Println("no chapters found (Cover/Title)")
	}
	return nil
}

// func (c *Chapter) saveIllustration() error {
// 	var link *url.URL
// 	for i := range c.ChapterBody {
// 		p := c.ChapterBody[i]
// 		if strings.Contains(p, "<img") {
// 			reader := strings.NewReader(p)
// 			doc, err := goquery.NewDocumentFromReader(reader)
// 			if err != nil {
// 				return err
// 			}
// 			// TODO fix this
// 			s := doc.Find("img")
// 			href, exists := s.Attr("img")
// 			if !exists {
// 				return fmt.Errorf("something went wrong with saveIllustration. (jqery and strings.Contains disagree)")
// 			}
// 			link, err = url.Parse(href)
// 			if err != nil {
// 				return err
// 			}
// 			filePath, err := saveHelper(link.String())
// 			if err != nil {
// 				return err
// 			}
// 			c.ChapterBody[i] = fmt.Sprintf("<p><img src=\"%s\" alt=\"\"/></p>", filePath)
// 			c.ChapterImages = append(c.ChapterImages, filePath)
// 		}
// 	}
// 	return nil
// }

func (v *Volume) saveCover() error {
	filePath, err := saveHelper(v.Chapters[0].bookmarkData.Cover)
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
	paragraphs := firstP.AddSelection(rest)

	paragraphs.EachWithBreak(func(i int, s *goquery.Selection) bool {
		// Check if the current paragraph contains the specific "Previous" link
		// You can check for the text or the href attribute
		link := s.Find("a")
		if link.AttrOr("data-type", "") == "fcn_chapter" {
			return false // Break the EachWithBreak loop
		}
		// find the hr element and break out
		if s.Is("hr") {
			return false
		}

		// find the img tag
		img := s.Find("img")
		if img.Length() > 0 {
			// 1. Extract the src attribute
			src, exists := img.Attr("src")
			if exists && src != "" {
				// 2. Pass to your helper (return new file path)
				filePath, err := saveHelper(src)
				if err == nil {
					// 3. Create the new string and replace the element
					c.ChapterImages = append(c.ChapterImages, filePath)
					newHTML := fmt.Sprintf("<p><img src=\"%s\" alt=\"\"/></p>", filePath)
					s.ReplaceWithHtml(newHTML)

					// fmt.Printf("Chapter: %s line: %d imgName: %v\n", c.ChapterTitle, i, filePath)
					results = append(results, newHTML)
					return true
				}
			}
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

func saveHelper(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("received empty URL in saveHelper")
	}

	pathBits := strings.Split(url, "/")
	imgName := pathBits[len(pathBits)-1]

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	if err != nil {
		return "", nil
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	filePath := filepath.Join("internal", "scraper", "images", imgName)
	filePath = strings.ReplaceAll(filePath, " ", "-")
	out, err := os.Create(filePath)
	if err != nil {
		return "", nil
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return filePath, err
}
