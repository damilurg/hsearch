package parser

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/comov/gilles_search_kg/structs"

	"github.com/PuerkitoBio/goquery"
)

// [NOT]: Тема-негативка. Только факты. Арендаторам и Арендодателям внимание!
const negativeTheme = 2477961

// LoadNewOffers - получает страницу по URL, перебирает объявления и возвращает
// список объявдений
func LoadNewOffers(url string) (map[uint64]string, error) {
	doc, err := GetDocumentByUrl(url)
	if err != nil {
		return nil, err
	}

	offers := make(map[uint64]string)

	// находим все топики на странице и проходися по ним, заполняя список
	// найденых объявлений
	doc.Find(".topic_title").Each(func(i int, s *goquery.Selection) {

		// получаем url объявления
		href, ok := s.Attr("href")
		if !ok {
			return
		}

		offerId, err := idFromHref(href)
		if err != nil {
			log.Println("Can't get Id from href with an error", err)
			return
		}

		offers[offerId] = href
	})

	delete(offers, negativeTheme)
	return offers, nil
}

// LoadOffersDetail - выгружает и парсит offers по href
func LoadOffersDetail(offersList map[uint64]string) []*structs.Offer {
	// TODO: это нужно сделать в горутинах
	offers := make([]*structs.Offer, 0)
	for id, href := range offersList {
		doc, err := GetDocumentByUrl(href)
		if err != nil {
			log.Printf("Can't load offer %s with an error %s\f", href, err)
			continue
		}
		offers = append(offers, structs.ParseNewOffer(href, id, doc))
	}
	return offers
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

// idFromHref - получение Id с URL предложения
func idFromHref(href string) (uint64, error) {
	urlPath, err := url.Parse(href)
	if err != nil {
		return 0, err
	}

	id := urlPath.Query().Get("showtopic")
	if id == "" {
		return 0, fmt.Errorf("id not contain in href")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	return uint64(idInt), nil
}
