package parser

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aastashov/house_search_assistant/structs"

	"github.com/PuerkitoBio/goquery"
)

// Parse - получает страницу по URL, перебирает объявления и возвращает
// список объявдений
func Parse(url string) ([]*structs.Offer, error) {
	doc, err := GetDocumentByUrl(url)
	if err != nil {
		return nil, err
	}

	var offers []*structs.Offer

	// находим все топики на странице и проходися по ним, заполняя список
	// найденых объявлений
	doc.Find(".topic_title").Each(func(i int, s *goquery.Selection) {

		// получаем url объявления
		href, ok := s.Attr("href")
		if !ok {
			return
		}

		// на основной странице почти нет никакой информации, по этому идем на
		// страницу и вытаскиваем больше информации об объявлении
		doc, err := GetDocumentByUrl(href)
		if err != nil {
			log.Println("[Parse.GetDocumentByUrl] error:", err)
			return
		}

		offers = append(offers, structs.ParseNewOffer(href, doc))
	})

	return offers, nil
}

// GetDocumentByUrl - получает страницу по http, читает и возвращет объект
// goquery.Document для парсинга
func GetDocumentByUrl(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("[GetDocumentByUrl.Get] error:", err)
		return nil, err
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Println("[GetDocumentByUrl.defer.Close] error:", err)
		}
	}()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return goquery.NewDocumentFromReader(res.Body)
}
