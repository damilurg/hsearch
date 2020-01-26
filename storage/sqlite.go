package storage

import (
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/aastashov/house_search_assistant/structs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

type (
	// Connector - структура для храниения и управления подключением к бд
	Connector struct {
		DB *sql.DB
	}
)

var regexContain = regexp.MustCompile(`UNIQUE constraint failed*`)

// New - возвращает коннектор для подключения к базе данных. Код не должен знать
// какая бд или какой драйвер используется для работы с базой.
func New() (*Connector, error) {
	db, err := sql.Open("sqlite3", "house_search_assistant.db?cache=shared")
	if err != nil {
		return nil, err
	}

	// Для БД sqlite максимальное кол-во соединений желательно иметь 1, так как
	// если писать в бд будут 2 потока, то файл может быть покрашен, что несет
	// потерб данных. Это не production проет, по этому нет смысла в PG ил MySql
	db.SetMaxOpenConns(1)

	return &Connector{DB: db}, nil
}

func (c *Connector) Migrate(path string) error {
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return err
	}
	err = goose.Run("up", c.DB, path)
	if err == goose.ErrNoNextVersion {
		return nil
	}

	return err
}

// WriteOffer - записывает Offer в базу вместе с картинками и вазвращает Id в
// структуру
func (c *Connector) WriteOffer(offer *structs.Offer) error {
	_, err := c.DB.Exec(`INSERT INTO offer (
		ex_id,
		url,
		topic,
		price,
		phone,
		room_numbers,
		body) VALUES (?, ?, ?, ?, ?, ?, ?);`,

		offer.ExId,
		offer.Url,
		offer.Topic,
		offer.Price,
		offer.Phone,
		offer.RoomNumber,
		offer.Body,
	)
	if err != nil && !regexContain.MatchString(err.Error()) {
		return err
	}

	offerId := uint64(0)
	err = c.DB.QueryRow("SELECT id FROM offer WHERE ex_id = ?", offer.ExId).Scan(&offerId)
	if err != nil {
		return err
	}
	offer.Id = offerId
	return c.writeImages(strconv.Itoa(int(offerId)), offer.Images)
}

// writeImages - так как картинки храняться в отдельной таблице, то пишем мы их
// отдельно
func (c *Connector) writeImages(offerId string, images []string) error {
	if len(images) <= 0 {
		return nil
	}

	params := make([]interface{}, 0)

	paramsPattern := ""
	sep := ""
	for _, image := range images {
		paramsPattern += sep + "(?, ?)"
		sep = ", "
		params = append(params, offerId, image)
	}

	query := "INSERT INTO image (offer_id, path) VALUES " + paramsPattern
	_, err := c.DB.Exec(query, params...)
	if err != nil && !regexContain.MatchString(err.Error()) {
		return err
	}

	return nil
}

// ReadUsersForOrder - достает пользователей для которых нужно сделать рассылку
func (c *Connector) ReadUsersForOrder(offer *structs.Offer) ([]*structs.User, error) {
	rows, err := c.DB.Query(`
	SELECT 
       username,
       chat
	FROM user
	LEFT JOIN user_to_offer uto on (user.id = uto.user_id and dislike = 0)
	WHERE enabled = 1;`)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("[ReadUsersForOrder.Close] error:", err)
		}
	}()

	users := make([]*structs.User, 0)
	for rows.Next() {
		user := new(structs.User)
		err := rows.Scan(
			&user.Username,
			&user.Chat,
		)
		if err != nil {
			log.Println("[ReadUsersForOrder.Scan] error:", err)
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *Connector) StartSearch(chat int64, username string) error {
	user := &structs.User{
		Username: username,
		Chat:     chat,
	}

	err := c.ReadUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.WriteUser(user)
		}
	}

	_, err = c.DB.Exec("UPDATE user SET enabled = 1 WHERE username = ?;", user.Username)
	return err
}

func (c *Connector) WriteUser(user *structs.User) error {
	_, err := c.DB.Exec(
		"INSERT INTO user (username, chat) VALUES (?, ?);",
		user.Username,
		user.Chat,
	)
	return err
}

func (c *Connector) ReadUser(user *structs.User) error {
	err := c.DB.QueryRow(
		"SELECT id, chat, enabled FROM user WHERE username = ?;",
		user.Username,
	).Scan(
		&user.Id,
		&user.Chat,
		&user.Enable,
	)
	return err
}

func (c *Connector) StopSearch(username string) error {
	user := &structs.User{Username: username}
	err := c.ReadUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	if user.Enable {
		_, err = c.DB.Exec("UPDATE user SET enabled = 0 WHERE username = ?;", user.Username)
	}

	return err
}

func (c *Connector) Bookmarks(username string) ([]*structs.Offer, int64, error) {
	user := &structs.User{Username: username}
	err := c.ReadUser(user)
	if err != nil {
		return nil, 0, err
	}

	rows, err := c.DB.Query(`SELECT
       id,
       ex_id,
       url,
       topic,
       price,
       phone,
       room_numbers,
       body
	FROM offer
	LEFT JOIN user_to_offer uto on (offer.id = uto.offer_id AND uto.user_id = ?)
	WHERE like = 1;
;`,
		user.Id,
	)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("[Bookmarks.Close] error:", err)
		}
	}()

	offers := make([]*structs.Offer, 0)
	for rows.Next() {
		offer := new(structs.Offer)
		err := rows.Scan(
			&offer.Id,
			&offer.ExId,
			&offer.Url,
			&offer.Topic,
			&offer.Price,
			&offer.Phone,
			&offer.RoomNumber,
			&offer.Body,
		)
		if err != nil {
			log.Println("[Bookmarks.Scan] error:", err)
			continue
		}
		offers = append(offers, offer)
	}

	return offers, user.Chat, nil
}

func (c *Connector) SaveMessage(msgId int, offerId uint64, chat int64) error {
	_, err := c.DB.Exec("INSERT INTO save_message (message_id, offer_id, chat) VALUES (?, ?, ?);",
		msgId,
		offerId,
		chat,
	)

	if err != nil && !regexContain.MatchString(err.Error()) {
		return nil
	}
	return err
}

func (c *Connector) Dislike(msgId int, user *structs.User) error {
	offerId := uint64(0)
	err := c.DB.QueryRow(
		"SELECT offer_id FROM save_message WHERE message_id = ? AND chat = ?;",
		msgId,
		user.Chat,
	).Scan(
		&offerId,
	)
	if err != nil {
		return err
	}

	err = c.ReadUser(user)
	if err != nil {
		return err
	}

	_, err = c.DB.Exec("INSERT INTO user_to_offer (user_id, offer_id, dislike) VALUES (?, ?, ?);", user.Id, offerId, 1)
	return err
}

func (c *Connector) ReadNextOffer(user *structs.User) (*structs.Offer, error) {
	err := c.ReadUser(user)
	if err != nil {
		return nil, nil
	}

	offer := new(structs.Offer)
	now := time.Now()
	// TODO: у меня голова бодит, этот sql не правельный даже без запуска понятно
	err = c.DB.QueryRow(`
	SELECT id,
		ex_id,
		url,
		topic,
		price,
		phone,
		room_numbers,
		body
	FROM offer
	LEFT JOIN user_to_offer uto on (offer.id = uto.offer_id AND uto.skip <= ? AND uto.dislike = 0 AND uto.user_id = ?)
	`,
		now.Unix(),
		user.Id,
	).Scan(
		&offer.Id,
		&offer.ExId,
		&offer.Url,
		&offer.Topic,
		&offer.Price,
		&offer.Phone,
		&offer.RoomNumber,
		&offer.Body,
	)
	return offer, err
}
