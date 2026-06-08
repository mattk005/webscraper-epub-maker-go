package main

import (
	"fmt"
	"os"
	"testing"
)

func getTestData() string {
	file, err := os.ReadFile("testdata/prologue.html")
	if err != nil {
		fmt.Println(err)
	}
	return string(file)
}

func TestGetTitle(t *testing.T) {
	inputHTML := getTestData()
	expected := "TNG Vol. 1 Prologue"

	// Call the function
	result, err := getTitle(inputHTML)
	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestGetBody(t *testing.T) {
	inputHTML := getTestData()
	// expected first and last string
	expected := []string{"<p>In the deepest area of the Gate to the Netherworld lay “THE NEW GATE”.</p>", "<p>The same instant, Shin lost his consciousness—</p>"}
	result, err := getBody(inputHTML)
	if len(result) != 165 {
		t.Errorf("expected 165 paragraphs, got %d", len(result))
	}
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result[0] != expected[0] || result[len(result)-1] != expected[1] {
		t.Errorf("expected %q, got %q", expected[0], result[0])
	}
}

func TestGetBookmarkData(t *testing.T) {
	inputHTML := getTestData()
	result, err := getBookmarkData(inputHTML)
	expected := BookmarkData{
		Thumbnail:  "https://shintranslations.com/wp-content/uploads/2017/05/Vol1Cover.jpg",
		Cover:      "https://shintranslations.com/wp-content/uploads/2017/05/Vol1Cover.jpg",
		Title:      "TNG Vol. 1 Prologue",
		StoryTitle: "TNG VOLUME 1: A Beginning And An End",
	}
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
