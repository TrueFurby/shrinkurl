package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	Url UrlRepo
}

func NewStorage(file string) *Storage {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}

	if _, err = db.Exec(urlSchema); err != nil {
		panic(err)
	}

	return &Storage{
		Url: &urlStorage{db},
	}
}

const urlSchema = `CREATE TABLE IF NOT EXISTS urls (
    id integer PRIMARY KEY,
    hash varchar(255) NOT NULL,
    url varchar(255) NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS duplicateUrl ON urls(hash);`

type urlStorage struct {
	*sql.DB
}

func (db *urlStorage) Update(u *Url) error {
	result, err := db.Exec("INSERT INTO urls (hash, url) VALUES (?, ?)", u.Hash, u.Url)
	if err != nil {
		return err
	}
	u.Id, err = result.LastInsertId()
	return err
}

func (db *urlStorage) Get(hash string) (u Url, err error) {
	row := db.QueryRow("SELECT id, hash, url FROM urls WHERE hash=?", hash)
	if err = row.Scan(&u.Id, &u.Hash, &u.Url); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	return u, nil
}

func (db *urlStorage) GetByUrl(url string) (u Url, err error) {
	row := db.QueryRow("SELECT id, hash, url FROM urls WHERE url=?", url)
	if err = row.Scan(&u.Id, &u.Hash, &u.Url); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	return u, nil
}

func (db *urlStorage) Remove(hash string) error {
	_, err := db.Exec("DELETE FROM urls WHERE hash=?", hash)
	return err
}
