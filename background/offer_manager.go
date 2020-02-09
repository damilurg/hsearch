// Запустив впервые этого менеджера, я получил высер из сообщений. Подумав как
// обойти не только этот высер, но и оповещать пользователя о новых
// предложениях, я решил что следует сделать 2 инструмента. 1-й будет просто
// искать и писать в бд новые offers, после чего будет трегерить 2-й
// инструмент который будет проходиться по всем пользователям и искать любое
// первое предложение за последнее N времени и высылать в чат
package background

/* TODO: для всего пакета
    1. Для менеджера лучше сделать отдельный логер в файл и более подробно
      логировать его работу, так как сейчас это немного черный ящик
    2. Парсер нужно переделать на горутины, так как мы тратим много времени на
      ожидание по сети, а могли бы в этот момент парсить данные и писать их в
      бд. Нужно рассмотреть reader, writer для реализации pipeline
    3. ...
*/

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
		WriteOffers(offer []*structs.Offer) (int, error)
		ReadUsersForMatching() ([]*structs.User, error)
		ReadNextOffer(user *structs.User) (*structs.Offer, error)
		CleanFromExistOrders(offers map[uint64]string) error
	}

	Bot interface {
		SendPreviewMessage(offer *structs.Offer, user *structs.User) error
	}

	OfferManager struct {
		st  Storage
		bot Bot
		url string
		cnf *configs.Config
	}
)


// кодичество предложений на одной странице в diesel
const offersOnPage = 39

// StartOfferManager - инициализирует менеджера
func StartOfferManager(url string, cnf *configs.Config, st Storage, bot Bot) *OfferManager {
	return &OfferManager{
		st:  st,
		bot: bot,
		url: url,
		cnf: cnf,
	}
}

// Start - запускает работу менеджера
func (m *OfferManager) Start() {
	m.parser()
}

// parser - парсит удаленнвые ресурсы, находит предложения и пишет в хранилище,
// после чего трегерит broker
func (m *OfferManager) parser() {
	// при первом запуске менеджера, он начнет первый парсинг через 2 секунды,
	// а после изменится на время из настроек (sleep = m.cnf.ManagerDelay)
	sleep := time.Second * 2

	log.Println("[offer_manager] Start Offer Manager")
	for {
		select {
		case <-time.After(sleep):
			sleep = m.cnf.ManagerDelay

			for i := 1; i <= m.cnf.MaxPage; i++ {
				target := fmt.Sprintf(m.url, i)
				log.Printf("[offer_manager] start parse %s\n", target)
				offersLinks, err := parser.LoadNewOffers(target)
				if err != nil {
					log.Println("[offer_manager.LoadNewOffers] error: ", err)
					continue
				}

				if len(offersLinks) == 0 {
					log.Println("[offer_manager] no offers for target", target)
					continue
				}

				err = m.st.CleanFromExistOrders(offersLinks)
				if err != nil {
					log.Println("[offer_manager.CleanFromExistOrders] error", err)
					continue
				}

				offers := parser.LoadOffersDetail(offersLinks)
				log.Printf("[offer_manager] Find %d new offers for targer %s", len(offers), target)

				newOffers, err := m.st.WriteOffers(offers)
				if err != nil {
					log.Println("[offer_manager.WriteOffer] error", err)
					continue
				}

				if newOffers < offersOnPage {
					m.cnf.MaxPage = 1
				}
			}

			// после того как мы нашли новые offer, мы начинаем рассылку
			m.broker()
		}
	}
}

// broker - вытягивает всех пользователй их бд и начинает для них рассылку. Если
// пользователь не увидел первое сообщение, то при следущем вызове от parser,
// broker должен отослать следующее если есть. Если нет, то пропускаем
// пользователя
func (m *OfferManager) broker() {
	users, err := m.st.ReadUsersForMatching()
	if err != nil {
		log.Println("[offer_manager.ReadUsersForOrder] error:", err)
		return
	}

	if len(users) <= 0 {
		log.Println("[offer_manager.users.len] no users")
		return
	}

	for _, user := range users {
		// TODO: каждый раз создавать горутину на пользователя, это жирно. Нужно
		//  сделать воркеров которые будут создаваться при старте, затем
		//  передавать им пользователей и они будут заниматься matching
		go m.matching(user)
	}
}

func (m *OfferManager) matching(user *structs.User) {
	log.Printf("[offer_manager] Start matching for user %s", user.Username)

	offer, err := m.st.ReadNextOffer(user)
	if err != nil {
		log.Printf("[offer_manager] Can't read offer for user %s with an error %s\n", user.Username, err)
		return
	}

	if offer == nil {
		log.Printf("[offer_manager] For user %s not new offers", user.Username)
		return
	}

	err = m.bot.SendPreviewMessage(offer, user)
	if err != nil {
		log.Printf("[offer_manager] Can't send message for user `%s` with an error %s\n", user.Username, err)
		return
	}

	log.Printf("[offer_manager] Successfully send offer %d for user %s\n", offer.Id, user.Username)
}
