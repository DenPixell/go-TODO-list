package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8) NOT NULL DEFAULT "",
    title   VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat  VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX idx_scheduler_date ON scheduler (date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	database, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err := database.Exec(schema); err != nil {
			return err
		}
	}

	db = database
	return nil
}

func Close() error {
	return db.Close()
}
