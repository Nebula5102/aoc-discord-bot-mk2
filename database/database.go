package database

import (
	"log"
	"database/sql"
	"time"

	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const fileDB = "database/competition.db"

type User struct {
	id int
	DiscordID string
	AocID string
	Score int
}

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
			discordID TEXT NOT NULL,
			dayNumber INTEGER NOT NULL,
			startTime DATETIME,
			endTime DATETIME
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

func UpdateID(discordID string, aocID string, score int) {
	db, err := sql.Open("sqlite3",fileDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`INSERT OR REPLACE INTO USER (discordID,aocID,score) VALUES (?,?,?);`,discordID,aocID,score)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully Updated:",discordID)
}

func Score(discordID string) int {
	db, err := sql.Open("sqlite3",fileDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	val, err := db.Query("SELECT score FROM USER WHERE discordID = ?;",discordID)
	if err != nil {
		log.Fatal(err)
	}
	var score int
	for val.Next() {
		if err := val.Scan(&score); err != nil {
			log.Fatal(err)
		}
	}

	return score 
}

func InsertDay(discordID string, t time.Time, day int) {
	db, err := sql.Open("sqlite3",fileDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO DAY (discordID,dayNumber,startTime) VALUES (?, ?, ?);`,discordID,day,t)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully inserted day:",discordID)
}


func PullCompetitionBoard(competitors *[]User) {
	db, err := sql.Open("sqlite3",fileDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT * FROM USER ORDER BY score;")
	if err != nil {
		log.Fatal(err)
	}
	
	for row.Next() {
		user := User{}
		if err := row.Scan(&user.id,&user.DiscordID,&user.AocID,&user.Score); err != nil {
			log.Fatal(err)
		}
		*competitors = append(*competitors, user)
	}
}

func UpdateDay(t time.Time, discordID string,day int) {
	db, err := sql.Open("sqlite3",fileDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`UPDATE DAY SET endTime = ? WHERE discordID = ? AND dayNumber = ?`,t,discordID,day)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Successfully Updated:",discordID)
}

func GrabTime(discordID string, day int) (time.Time, time.Time) {
	db, err := sql.Open("sqlite3",fileDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var start time.Time
	var end time.Time
	row, err := db.Query(
		`SELECT startTime 
		 FROM DAY 
		 WHERE discordID = ? AND dayNumber = ?
	`,discordID,day)
	for row.Next() {
		if err := row.Scan(&start); err != nil {
			log.Fatal(err)
		}
	}
	row, err = db.Query(
		`SELECT endTime
		 FROM DAY 
		 WHERE discordID = ? AND dayNumber = ?
	`,discordID,day)
	for row.Next() {
		if err := row.Scan(&end); err != nil {
			log.Fatal(err)
		}
	}
	log.Println(start,end)
	return start, end
}

func UpdateScore(discordID string, st time.Time, et time.Time) {
	db, err := sql.Open("sqlite3",fileDB)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := et.Sub(st).Minutes()
	var points int
	switch {
	case elapsed < 30:
		points = 10
	case elapsed < 60:
		points = 8
	case elapsed < 120:
		points = 6
	case elapsed < 240:
		points = 4
	case elapsed < 480:
		points = 2
	case elapsed >= 480:
		points = 1
	} 
	total := Score(discordID) + points

	_, err = db.Exec(`UPDATE USER SET score = ? WHERE discordID = ?`,total,discordID)
	if err != nil {
		log.Fatal(err)
	}
}
