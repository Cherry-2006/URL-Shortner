package database

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open(
		"sqlite",
		"file:app.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)",
	)
	if err != nil {
		log.Fatal(err)
	}
	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS counter (count INTEGER NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}
	var n int
	err = DB.QueryRow("SELECT COUNT(*) FROM counter").Scan(&n)
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		_, err = DB.Exec("INSERT INTO counter (count) VALUES (0)")
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS urlTable (LongUrl TEXT, ShortUrl TEXT UNIQUE, creationTime DATETIME)")
	if err != nil {
		log.Fatal(err)
	}
}

func GetNextCounter() uint64 {
	var count uint64

	err := DB.QueryRow(`
    UPDATE counter
    SET count = count + 1
    RETURNING count
`).Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	return count
}

func InsertURL(longUrl string, shortUrl string) error {
	_, err := DB.Exec("INSERT INTO urlTable VALUES (?, ?, ?)", longUrl, shortUrl, time.Now())
	return err
}

func GetLongURL(shortUrl string) (string, error) {
	var longUrl string
	err := DB.QueryRow(
		"SELECT LongUrl FROM urlTable WHERE ShortUrl=?",
		shortUrl,
	).Scan(&longUrl)
	return longUrl, err
}
