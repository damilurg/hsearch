package background

import (
	"log"
	"time"
)

func (m *Manager) garbage() {
	expireDate := time.Now().AddDate(0, 0, m.cnf.ExpireDays*-1).Unix()

	log.Printf("[garbage] StartGarbageCollector Manager\n")
	err := m.st.CleanExpiredOffers(expireDate)
	if err != nil {
		log.Printf("[garbage.CleanExpiredOffers] Error: %s\n", err)
	}

	err = m.st.CleanExpiredImages(expireDate)
	if err != nil {
		log.Printf("[garbage.CleanExpiredImages] Error: %s\n", err)
	}

	err = m.st.CleanExpiredAnswers(expireDate)
	if err != nil {
		log.Printf("[garbage.CleanExpiredAnswers] Error: %s\n", err)
	}

	err = m.st.CleanExpiredTGMessages(expireDate)
	if err != nil {
		log.Printf("[garbage.CleanExpiredTGMessages] Error: %s\n", err)
	}
}
