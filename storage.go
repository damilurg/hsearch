package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type (
	/*
		CREATE TABLE IF NOT EXIST user (
		    name VARCHAR(80)  DEFAULT '',
		    chat VARCHAR(80)  DEFAULT ''
		);
	*/
	User struct {
		Name string
		Chat string
	}

	/*
		CREATE TABLE IF NOT EXIST order (
		    link  VARCHAR(100)  DEFAULT '',
		    title VARCHAR(80)   DEFAULT '',
			price VARCHAR(50)   DEFAULT ''
		);
	*/
	Order struct {
		Link  string
		Title string
		Price string
	}
)

func New() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "realtor_bot.db?cache=shared")
	if db != nil {
		db.SetMaxOpenConns(1)
	}
	return db, err
}
