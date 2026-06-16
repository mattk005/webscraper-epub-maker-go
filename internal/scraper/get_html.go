package scraper

import (
	"fmt"
	"io"
	"net/http"
)

func getHTML(rawURL string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", rawURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// fmt.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	foo := http.DetectContentType(body)
	if foo != "text/html; charset=utf-8" {
		fmt.Println("Response Type != HTML")
		return "", fmt.Errorf("response Type != HTML")
	}
	return string(body), nil
}
