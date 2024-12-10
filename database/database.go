package database

import (
	"log"
	"database/sql"

	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const fileDB = "database/competition.db"

func InitTables() {
	db, err := sqlite3.Open(fileDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		log.Fatal(err)
	}

	//User Table
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS USER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			discordID TEXT NOT NULL UNIQUE,
			aocID TEXT NOT NULL UNIQUE,
			score INTEGER NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	//Day Table
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS DAY (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userID INTEGER NOT NULL,
			dayNumber INTEGER NOT NULL,
			startTime DATETIME,
			endTime DATETIME,
			FOREIGN KEY(userID) REFERENCES USER(id),
			UNIQUE(userID,dayNumber)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database Tables exist, or successfully created")
}

func UserSignup(discordID string, aocID string) {
	db, err := sql.Open("sqlite3",fileDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`INSERT OR IGNORE INTO USER (discordID,aocID,score) VALUES (?, ?, ?);`,discordID,aocID,0)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully inserted:",discordID)
}

func UpdateID(discordID string, aocID string) {
	db, err := sql.Open("sqlite3",fileDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`INSERT OR REPLACE INTO USER (discordID,aocID) VALUES (?,?);`,discordID,aocID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully Updated:",discordID)
}

/*
func PullCompetitionBoard() {
	db, err := sqlite3.Open(fileDB)
	if err != nil {
		log.Fatal(err)
	}
}
*/
