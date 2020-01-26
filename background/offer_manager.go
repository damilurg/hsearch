package background

import (
	"fmt"
	"log"
	"time"

	"github.com/aastashov/house_search_assistant/configs"
	"github.com/aastashov/house_search_assistant/parser"
	"github.com/aastashov/house_search_assistant/structs"
)

type (
	Storage interface {
		WriteOffer(offer *structs.Offer) error
		ReadUsersForOrder(offer *structs.Offer) ([]*structs.User, error)
	}

	Bot interface {
		SendPreviewMessage(offer *structs.Offer, user *structs.User) error
	}

	OfferManager struct {
		st  Storage
		bot Bot
	}
)


// TODO: это жесткий костыль, который нужно переделать
var endSedToUsers = false

func StartOfferManager(url string, cnf *configs.Config, st Storage, bot Bot) {
	endSedToUsers = true
	ofm := OfferManager{
		st:  st,
		bot: bot,
	}

	log.Println("[StartOfferManager]")
	for {
		select {
		case <-time.After(cnf.ManagerDelayMin * time.Minute):
			if !endSedToUsers {
				continue
			}

			for i := 1; i <= cnf.MaxPage; i++ {
				offers, err := parser.Parse(url)
				if err != nil {
					log.Println("[StartOfferManager.Parse] error: ", err)
					continue
				}
				ofm.writeOfferAndSendToUser(offers)
			}
			endSedToUsers = true
		}
	}
}

func (m *OfferManager) writeOfferAndSendToUser(offers []*structs.Offer) {
	for _, offer := range offers {
		err := m.st.WriteOffer(offer)
		if err != nil {
			fmt.Println("[writeOfferAndSendToUser.WriteOffer] error:", err)
			continue
		}

		users, err := m.st.ReadUsersForOrder(offer)
		if err != nil {
			fmt.Println("[writeOfferAndSendToUser.ReadUsersForOrder] error:", err)
			continue
		}

		if len(users) <= 0 {
			log.Println("[writeOfferAndSendToUser.users.len] no users")
			continue
		}

		for _, user := range users {
			err = m.bot.SendPreviewMessage(offer, user)
			if err != nil {
				fmt.Println("[writeOfferAndSendToUser.SendToUser] error:", err)
				continue
			}
		}
	}
}
