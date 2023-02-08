package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int64
	Channels string
	Time     int16
	Location int16
}

var DatabaseClient *sql.DB

func InitDatabase(config ServerConfig) {
	MySqlClient, err := sql.Open("mysql", config.SqlUser+":"+config.SqlPass+"@/postowl")
	if err != nil {
		log.Fatal(err)
	}

	MySqlClient.SetMaxOpenConns(config.MaxSqlConns)
	MySqlClient.SetMaxIdleConns(config.MaxSqlIdleConns)

	err = MySqlClient.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DatabaseClient = MySqlClient
}

func (user *User) Create() {
	_, err := DatabaseClient.Exec(
		"INSERT postowl.users(id, channels, time, location) VALUES(?, ?, ?, ?)",
		user.ID, user.Channels, user.Location, user.Time)
	if err != nil {
		log.Fatal(err)
	}
}

func (user *User) Get() bool {
	result := DatabaseClient.QueryRow("SELECT * FROM postowl.users WHERE id=?", user.ID)
	if result.Err() != nil {
		log.Fatal(result.Err())
	}

	err := result.Scan(&user.ID, &user.Channels, &user.Time, &user.Location)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return true
}

func (user *User) Update() {
	_, err := DatabaseClient.Exec(
		"UPDATE postowl.users SET time=?, channels=?, location=? WHERE id=?",
		user.Time, user.Channels, user.Location, user.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func DatabaseForScheduler(time int16) []int64 {
	rows, err := DatabaseClient.Query("SELECT id FROM postowl.users WHERE time=?", time)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var id int64
	var ids []int64
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Println(err)
		}

		ids = append(ids, id)
	}

	return ids
}
