package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

type (
	// Offer - хранит все объявления нашедшие на diesel
	Offer struct {
		Id     uint64
		Topic  string
		Body   string
		Images []string
		Price  string
		ExId   uint64
	}

	// UsersToOffers - это ManyToMany для хранения реакции пользователя на
	// объявдение
	UsersToOffers struct {
		UserId  uint64
		OfferId uint64
		Like    bool
		Dislike bool
	}

	// User - telegram пользователь
	User struct {
		Id      uint64
		Account string
		Chat    string
	}

	// Connector - структура для храниения и управления подключением к бд
	Connector struct {
		DB          *sqlx.DB
		createOffer *sqlx.Stmt
	}
)

// New - возвращает коннектор для подключения к базе данных. Код не должен знать
// какая бд или какой драйвер используется для работы с базой.
func New() (*Connector, error) {
	db, err := sqlx.Connect("sqlite3", "house_search_assistant.db?cache=shared")
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
	err = goose.Run("up", c.DB.DB, path)
	if err == goose.ErrNoNextVersion {
		return nil
	}

	return err
}
