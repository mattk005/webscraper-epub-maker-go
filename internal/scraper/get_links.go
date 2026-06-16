package scraper

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	var links []string

	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return links, err
	}
	linksSelection := doc.Find("a[href]")
	linksSelection.Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}
		u, err := url.Parse(href)
		if err != nil {
			// fmt.Printf("couldn't prase href %q: %v\n", href, err)
			return
		}

		resolved := baseURL.ResolveReference(u)
		// fmt.Printf("DEBUG: Appending %s to links slice\n", resolved.String())
		links = append(links, resolved.String())
	})

	return links, nil
}
