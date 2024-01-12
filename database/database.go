package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

var dbConfig = mysql.Config{User: "root", Passwd: "anurag", Net: "tcp", Addr: "0.0.0.0:3306", DBName: "urldb"}

type Database struct {
	conn *sql.DB
}

func New() Database {
	conn, err := sql.Open("mysql", dbConfig.FormatDSN())
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)
	if err != nil {
		panic("DB connection failure: " + err.Error())
	}
	if err := conn.Ping(); err != nil {
		panic("DB ping failure: " + err.Error())
	}
	log.Println("DB connection successful")
	return Database{conn}
}

func(db *Database) Close() error {
	return db.conn.Close()
}

func(db *Database) ShortUrlExists(shortUrl string) (bool, error) {
	row := db.conn.QueryRow("select count(*) from urls where shorturl = ?", shortUrl)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func(db *Database) LongUrlExists(longUrl string) (bool, error) {
	row := db.conn.QueryRow("select count(*) from urls where longurl = ?", longUrl)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *Database) GetLongUrl(shortUrl string) (string, error) {
	row := db.conn.QueryRow("select longurl from urls where shorturl = ?", shortUrl)
	var longUrl string
	if err := row.Scan(&longUrl); err != nil {
		return "", err
	}
	return longUrl, nil
}

func(db *Database) GetShortUrl(longUrl string) (string, error) {
	row := db.conn.QueryRow("select shorturl from urls where longurl = ?", longUrl)
	var shortUrl string
	if err := row.Scan(&shortUrl); err != nil {
		return "", err
	}
	return shortUrl, nil
}

func(db *Database) InsertUrl(shortUrl, longUrl string) error {
	_, err := db.conn.Exec("insert into urls values (?, ?)", shortUrl, longUrl)
	return err
}