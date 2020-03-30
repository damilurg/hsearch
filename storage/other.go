package storage

import (
	"log"
	"time"

	"github.com/comov/hsearch/structs"
)

// Feedback - write feedback from user
func (c *Connector) Feedback(chat int64, username, body string) error {
	_, err := c.DB.Exec(
		"INSERT INTO feedback (created, chat, username, body) VALUES (?, ?, ?, ?);",
		time.Now().Unix(),
		chat,
		username,
		body,
	)
	return err
}

// SaveMessage - when we send user or group offer, description or photos, we
//  save this message for subsequent removal from chat, if need.
func (c *Connector) SaveMessage(msgId int, offerId uint64, chat int64, kind string) error {
	_, err := c.DB.Exec(
		`INSERT INTO tg_messages (message_id, offer_id, kind, chat, created)
 				VALUES (?, ?, ?, ?, ?);`,
		msgId,
		offerId,
		kind,
		chat,
		time.Now().Unix(),
	)

	return err
}

// ReadChat - return chat with user or with group if exist or return error.
func (c *Connector) ReadChat(id int64) (*structs.Chat, error) {
	chat := &structs.Chat{}
	err := c.DB.QueryRow(
		`SELECT id, username, title, enable, c_type
				FROM chat
				WHERE id = ?;`,
		id,
	).Scan(
		&chat.Id,
		&chat.Username,
		&chat.Title,
		&chat.Enable,
		&chat.Type,
	)
	return chat, err
}


// ReadChatsForMatching - read all the chats for which the mailing list has
//  to be done
func (c *Connector) ReadChatsForMatching() ([]*structs.Chat, error) {
	rows, err := c.DB.Query(`
	SELECT DISTINCT
		c.id,
		c.username,
		c.title
	FROM chat c
	LEFT JOIN answer uto on (c.id = uto.chat and uto.dislike = 0)
	WHERE c.enable = 1;`)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("[ReadChatsForMatching.Close] error:", err)
		}
	}()

	chats := make([]*structs.Chat, 0)
	for rows.Next() {
		chat := new(structs.Chat)
		err := rows.Scan(
			&chat.Id,
			&chat.Username,
			&chat.Title,
		)
		if err != nil {
			log.Println("[ReadChatsForMatching.Scan] error:", err)
			continue
		}

		chats = append(chats, chat)
	}

	return chats, nil
}
