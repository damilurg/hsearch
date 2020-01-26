package structs

import (
	"log"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type (
	// Offer - хранит все объявления c diesel
	Offer struct {
		Id         uint64
		ExId       uint64
		Url        string
		Topic      string
		Price      string
		Phone      string
		RoomNumber string
		Body       string
		Images     []string
		doc        *goquery.Document
	}
)

// ParseNewOffer - заполняет структуру объявления
func ParseNewOffer(href string, doc *goquery.Document) *Offer {
	offer := &Offer{
		Url: href,
		doc: doc,
	}

	offer.parseId(href)
	offer.parseTitle()
	offer.parsePrice()
	offer.parsePhone()
	offer.parseRoomNumber()
	offer.parseBody()
	offer.parseImages()
	return offer
}

// parseId - достает Id из URL
func (o *Offer) parseId(href string) uint64 {
	urlPath, err := url.Parse(href)
	if err != nil {
		log.Println("[parseId.Parse] error:", err)
		return 0
	}

	id := urlPath.Query().Get("showtopic")
	if id == "" {
		log.Println("[parseId.Get] id is empty")
		return 0
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println("[parseId.Atoi] error:", err)
		return 0
	}

	o.ExId = uint64(idInt)
	return o.ExId
}

// parseTitle - находит заголовок
func (o *Offer) parseTitle() string {
	o.Topic = o.doc.Find(".ipsType_pagetitle").Text()
	return o.Topic
}

// parsePrice - находит цену их баджика
func (o *Offer) parsePrice() string {
	o.Price = o.doc.Find("span.field-value.badge.badge-green").Text()
	return o.Price
}

// parsePhone - находит номер телефона из настроек обхявления
func (o *Offer) parsePhone() string {
	o.Phone = o.doc.Find(".custom-field.md-phone > span.field-value").Text()
	return o.Phone
}

// parseRoomNumber - находит количество комнат из настроек объявления
func (o *Offer) parseRoomNumber() string {
	roomNumberNodes := o.doc.Find("span:contains('Количество комнат')").Parent().Children().Nodes
	if len(roomNumberNodes) > 1 {
		o.RoomNumber = goquery.NewDocumentFromNode(roomNumberNodes[1]).Text()
	}
	return o.RoomNumber
}

// parseBody - находит тело объявления
func (o *Offer) parseBody() string {
	// todo: нужно почистить от html тегов
	o.Body = o.doc.Find(".post.entry-content").Text()
	return o.Body
}

// parseImages - находит все прикрепленные файлы в объявлении
func (o *Offer) parseImages() []string {
	o.doc.Find(".attach").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("src")
		if ok {
			o.Images = append(o.Images, href)
		}
	})
	return o.Images
}
