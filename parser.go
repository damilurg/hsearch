package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Parse(callBack func(string)) {
	res, err := http.Get(BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".topic_title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		text := strings.TrimSpace(s.Text())
		href, ok := s.Attr("href")
		if !ok {
			return
		}

		callBack(fmt.Sprintf("Review %d: `%s` - %s\n", i, text, href))
	})
}
