package main

import (
	"log"

	"wp-to-epub/internal/scraper"
)

func main() {
	volumes := scraper.ScrapeAndParse()
	err := makeEpub(volumes)
	if err != nil {
		log.Fatal(err)
	}
}
