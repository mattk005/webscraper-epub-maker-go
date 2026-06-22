package scraper

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func getTestData() string {
	file, err := os.ReadFile("testdata/chap3.html")
	if err != nil {
		fmt.Println(err)
	}
	return string(file)
}

func testSaveIllustration(t *testing.T) {
	inputHTML := getTestData()
	// chapterHTML, err := getHTML(chapterURL)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	reader := strings.NewReader(inputHTML)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	chap := Chapter{}
	chap.getTitle(doc)
	chap.getBody(doc)
	chap.getBookmarkData(doc)

	err = chap.saveIllustration()
	// expected := ""
	// if err != nil {
	// 	t.Fatalf("expected no error, got %v", err)
	// }
	// if result != expected {
	// 	t.Errorf("expected %q, got %q", expected, result)
	// }
}
