package main

import (
	"log"

	"github.com/comov/hsearch/parser"
)

func main() {
	//var site = parser.DieselSite()
	var site = parser.HouseSite()
	//var site = parser.LalafoSite()

	//doc, err := parser.GetDocumentByUrl(site.Url())
	//if err != nil {
	//	log.Fatalln(err)
	//}

	offersLinks, err := parser.FindOffersLinksOnSite(site)
	if err != nil {
		log.Fatalln(err)
	}

	for id, link := range offersLinks {
		log.Println(id, link)
	}
}
