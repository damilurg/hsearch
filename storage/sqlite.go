package storage

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/comov/gilles_search_kg/configs"
	"github.com/comov/gilles_search_kg/structs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

type (
	// Connector - структура для храниения и управления подключением к бд
	Connector struct {
		DB              *sql.DB
		skipTime        time.Duration
		freshOffersTime time.Duration
	}
)

var regexContain = regexp.MustCompile(`UNIQUE constraint failed*`)

// New - возвращает коннектор для подключения к базе данных. Код не должен знать
// какая бд или какой драйвер используется для работы с базой.
func New(cnf *configs.Config) (*Connector, error) {
	db, err := sql.Open("sqlite3", "gilles_search_kg.db?cache=shared")
	if err != nil {
		return nil, err
	}

	// Для БД sqlite максимальное кол-во соединений желательно иметь 1, так как
	// если писать в бд будут 2 потока, то файл может быть покрашен, что несет
	// потерб данных. Это не production проет, по этому нет смысла в PG ил MySql
	db.SetMaxOpenConns(1)

	return &Connector{
		DB:              db,
		skipTime:        cnf.SkipTime,
		freshOffersTime: cnf.FreshOffers,
	}, nil
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
		id,
		url,
		topic,
		price,
		phone,
		room_numbers,
    	body,
		images,
		created) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		offer.Id,
		offer.Url,
		offer.Topic,
		offer.Price,
		offer.Phone,
		offer.Rooms,
		offer.Body,
		offer.Images,
		time.Now().Unix(),
	)
	if err != nil && !regexContain.MatchString(err.Error()) {
		return err
	}
	return c.writeImages(strconv.Itoa(int(offer.Id)), offer.ImagesList)
}

// WriteOffers - пишет пачку из offers вместе с картинками в бд
func (c *Connector) WriteOffers(offers []*structs.Offer) (int, error) {
	newOffersCount := 0
	// TODO: как видно, сейчас это сделано через простой цикл, но лучше
	//  предоставить это самому хранилищу. Сделать bulk insert, затем запросить
	//  Id по ExtId и записать картины. Не было времени сделать это сразу
	for i := range offers {
		offer := offers[i]
		err := c.WriteOffer(offer)
		if err != nil {
			return newOffersCount, err
		}

		newOffersCount += 1
	}
	return newOffersCount, nil
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

// CleanFromExistOrders - очищает map от offers которые уже есть в базе
func (c *Connector) CleanFromExistOrders(offers map[uint64]string) error {
	params := make([]interface{}, 0)

	paramsPattern := ""
	sep := ""
	for id := range offers {
		paramsPattern += sep + "?"
		sep = ", "
		params = append(params, id)
	}

	query := fmt.Sprintf("SELECT id FROM offer WHERE id IN (%s)", paramsPattern)
	rows, err := c.DB.Query(query, params...)
	if err != nil {
		return err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("[CleanFromExistOrders.Close] error:", err)
		}
	}()

	for rows.Next() {
		exId := uint64(0)
		err := rows.Scan(&exId)
		if err != nil {
			log.Println("[CleanFromExistOrders.Scan] error:", err)
			continue
		}

		delete(offers, exId)
	}

	return nil
}

// ReadUsersForOrder - достает пользователей для которых нужно сделать рассылку
func (c *Connector) ReadUsersForMatching() ([]*structs.User, error) {
	rows, err := c.DB.Query(`
	SELECT DISTINCT
       u.username,
       u.chat
	FROM user u
	LEFT JOIN answer uto on (u.chat = uto.chat and uto.dislike = 0)
	WHERE u.enable = 1;`)
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
			err = c.WriteUser(user)
			if err != nil {
				return err
			}
		}
	}

	_, err = c.DB.Exec("UPDATE user SET enable = 1 WHERE username = ?;", user.Username)
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
		"SELECT chat, enable FROM user WHERE username = ?;",
		user.Username,
	).Scan(
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
		_, err = c.DB.Exec("UPDATE user SET enable = 0 WHERE username = ?;", user.Username)
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
       url,
       topic,
       price,
       phone,
       room_numbers,
       images
	FROM offer
	LEFT JOIN answer uto on (offer.id = uto.offer_id AND uto.chat = ?)
	WHERE like = 1;`,
		user.Chat,
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
			&offer.Url,
			&offer.Topic,
			&offer.Price,
			&offer.Phone,
			&offer.Rooms,
			&offer.Images,
		)
		if err != nil {
			return nil, 0, err
		}
		offers = append(offers, offer)
	}

	return offers, user.Chat, nil
}

func (c *Connector) SaveMessage(msgId int, offerId uint64, chat int64) error {
	_, err := c.DB.Exec("INSERT INTO tg_messages (message_id, offer_id, chat, created) VALUES (?, ?, ?, ?);",
		msgId,
		offerId,
		chat,
		time.Now().Unix(),
	)

	if err != nil && !regexContain.MatchString(err.Error()) {
		return nil
	}
	return err
}

func (c *Connector) Dislike(msgId int, user *structs.User) error {
	offerId := uint64(0)
	err := c.DB.QueryRow(
		"SELECT offer_id FROM tg_messages WHERE message_id = ? AND chat = ?;",
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

	_, err = c.DB.Exec(
		"INSERT INTO answer (chat, offer_id, dislike, created) VALUES (?, ?, ?, ?);",
		user.Chat,
		offerId,
		1,
		time.Now().Unix(),
	)
	return err
}

func (c *Connector) Skip(msgId int, user *structs.User) error {
	offerId := uint64(0)
	err := c.DB.QueryRow(
		"SELECT offer_id FROM tg_messages WHERE message_id = ? AND chat = ?;",
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

	skipTime := time.Now().Add(c.skipTime).Unix()
	_, err = c.DB.Exec(
		"INSERT INTO answer (chat, offer_id, skip, created) VALUES (?, ?, ?, ?);",
		user.Chat,
		offerId,
		skipTime,
		time.Now().Unix(),
	)
	return err
}

func (c *Connector) ReadNextOffer(user *structs.User) (*structs.Offer, error) {
	err := c.ReadUser(user)
	if err != nil {
		return nil, nil
	}

	offer := new(structs.Offer)
	now := time.Now()

	err = c.DB.QueryRow(`
	SELECT DISTINCT
	   id,
       url,
       topic,
       price,
       phone,
       room_numbers,
	   images,
	   body
	FROM offer of
	LEFT JOIN answer u on (of.id = u.offer_id AND u.chat = ?)
	LEFT JOIN tg_messages sm on (of.id = sm.offer_id AND sm.chat = ?)
	WHERE of.created >= ?
	  AND (u.dislike = 0 OR u.dislike IS NULL)
	  AND sm.created IS NULL
	ORDER BY of.created;
	`,
		user.Chat,
		user.Chat,
		now.Add(-c.freshOffersTime).Unix(),
	).Scan(
		&offer.Id,
		&offer.Url,
		&offer.Topic,
		&offer.Price,
		&offer.Phone,
		&offer.Rooms,
		&offer.Images,
		&offer.Body,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return offer, nil
}

func (c *Connector) ReadOfferDescription(msgId int, user *structs.User) (string, error) {
	offerId := uint64(0)
	err := c.DB.QueryRow(
		"SELECT offer_id FROM tg_messages WHERE message_id = ? AND chat = ?;",
		msgId,
		user.Chat,
	).Scan(
		&offerId,
	)
	if err != nil {
		return "", err
	}

	description := ""
	err = c.DB.QueryRow(`SELECT body FROM offer of WHERE of.id = ?;`,
		offerId,
	).Scan(
		&description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return "Предложение не найдено, возможно было удалено", nil
		}
		return "", err
	}

	return description, nil
}

func (c *Connector) ReadOfferImages(msgId int, user *structs.User) ([]string, error) {
	offerId := uint64(0)
	images := make([]string, 0)

	err := c.DB.QueryRow(
		"SELECT offer_id FROM tg_messages WHERE message_id = ? AND chat = ?;",
		msgId,
		user.Chat,
	).Scan(
		&offerId,
	)
	if err != nil {
		return images, err
	}

	rows, err := c.DB.Query(`SELECT path FROM image im WHERE im.offer_id = ?;`, offerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return images, nil
		}
		return images, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("[ReadOfferImages] error:", err)
		}
	}()

	for rows.Next() {
		image := ""
		err := rows.Scan(
			&image,
		)
		if err != nil {
			log.Println("[ReadOfferImages.Scan] error:", err)
			continue
		}
		images = append(images, image)
	}

	return images, nil
}

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
