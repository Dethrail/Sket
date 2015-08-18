package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var (
	db *sql.DB
)

func OpenDBConnection() {
	var err error
	db, err = sql.Open("postgres", "user=Dethrail password=+ dbname=Tamagotchi sslmode=disable") // sslmode=verify-full // this will require ssl
	if err != nil {
		log.Fatal(err)
	}
}

func SelectUsers() (s string) {

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int
		var name string
		var password string

		err = rows.Scan(&id, &name, &password)
		if err != nil {
			fmt.Printf("rows.Scan error: %v\n", err)
		}

		fmt.Printf("name: %v password: %v \n", name, password)
		return password
	}
	return ""
}

func SelectUsersByPassword(pwd string) {

	rows, err := db.Query("SELECT * FROM users WHERE password='" + pwd + "'")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int
		var name string
		var password string

		err = rows.Scan(&id, &name, &password)
		if err != nil {
			fmt.Printf("rows.Scan error: %v\n", err)
		}
		fmt.Printf("with pwd name: %v password: %v \n", name, password)
	}
}
