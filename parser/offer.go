package parser

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/comov/hsearch/structs"

	"github.com/PuerkitoBio/goquery"
)

var (
	priceRegex    = regexp.MustCompile(`\d+`)
	currencyRegex = regexp.MustCompile(`[a-zA-Zа-яА-Я]+`)
)

// ParseNewOffer - parse html and fills the offer with valid values
func ParseNewOffer(href string, exId uint64, doc *goquery.Document) *structs.Offer {
	fullPrice, price, currency := parsePrice(doc)
	images := parseImages(doc)
	return &structs.Offer{
		Id:         exId,
		Url:        href,
		Topic:      parseTitle(doc),
		FullPrice:  fullPrice,
		Price:      price,
		Currency:   currency,
		Phone:      parsePhone(doc),
		Rooms:      parseRoomNumber(doc),
		Body:       parseBody(doc),
		Images:     len(images),
		ImagesList: images,
	}
}

// parseTitle - find topic title
func parseTitle(doc *goquery.Document) string {
	return doc.Find(".ipsType_pagetitle").Text()
}

// parsePrice - find price from badge
func parsePrice(doc *goquery.Document) (string, int, string) {
	fullPrice := doc.Find("span.field-value.badge.badge-green").Text()
	price := 0
	currency := ""

	pInt := priceRegex.FindAllString(fullPrice, -1)
	if len(pInt) == 1 {
		p, err := strconv.Atoi(pInt[0])
		if err != nil {
			log.Printf("[parsePrice] %s with an error: %s", fullPrice, err)
		}
		price = p
	}

	pCurrency := currencyRegex.FindAllString(fullPrice, -1)
	if len(pCurrency) == 1 {
		currency = strings.ToLower(pCurrency[0])
	}

	return fullPrice, price, currency
}

// parsePhone - find phone number from badge
func parsePhone(doc *goquery.Document) string {
	return doc.Find(".custom-field.md-phone > span.field-value").Text()
}

// parseRoomNumber - find number of rooms from badge
func parseRoomNumber(doc *goquery.Document) string {
	roomNumberNodes := doc.Find("span:contains('Количество комнат')").Parent().Children().Nodes
	rooms := ""
	if len(roomNumberNodes) > 1 {
		rooms = goquery.NewDocumentFromNode(roomNumberNodes[1]).Text()
	}
	return rooms
}

// parseBody - find offer body in page
func parseBody(doc *goquery.Document) string {
	messages := doc.Find(".post.entry-content").Nodes
	body := ""
	if len(messages) != 0 {
		body := goquery.NewDocumentFromNode(messages[0]).Text()
		reg := regexp.MustCompile(`Сообщение отредактировал.*`)
		body = reg.ReplaceAllString(body, "${1}")
		body = strings.Replace(body, "Прикрепленные изображения", "", 1)
		body = strings.Replace(body, "  ", "", 1)
		body = strings.TrimSpace(body)
	}
	return body
}

// parseImages - file all attachment in offer
func parseImages(doc *goquery.Document) []string {
	images := make([]string, 0)
	doc.Find(".attach").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("src")
		if ok {
			images = append(images, href)
		}
	})
	return images
}
