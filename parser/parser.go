package parser

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// TestSelector - отладочная функция нужная для тестирования отдельных
// select-оров. Я обычно пишу ее в main, a после указываю return и не запускаю
// полностью код. Если вносищь изменения, то потом затри. Не делай лишних doff-в
func TestSelector(url string) {
	doc, err := getDocumentByUrl(url)
	if err != nil {
		log.Println("[Debug.getDocumentByUrl] error: ", err)
		return
	}

	// Для примера использую это так:
	// title := doc.Find(".ipsType_pagetitle").Text()
	// fmt.Println("Title: ", title)
	_ = doc
}

// Parse - получает страницу по URL, перебирает объявления и возвращает
// список объявдений
func Parse(url string) ([]*Offer, error) {
	doc, err := getDocumentByUrl(url)
	if err != nil {
		return nil, err
	}

	var offers []*Offer

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
		doc, err := getDocumentByUrl(href)
		if err != nil {
			log.Println("[Parse.getDocumentByUrl] error: ", err)
			return
		}

		offers = append(offers, ParseNewOffer(href, doc))
	})

	return offers, nil
}

// getDocumentByUrl - получает страницу по http, читает и возвращет объект
// goquery.Document для парсинга
func getDocumentByUrl(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("[getDocumentByUrl.Get] error: ", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return goquery.NewDocumentFromReader(res.Body)
}
