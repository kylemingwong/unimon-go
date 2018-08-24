package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1q2w3e$R%T"
	dbname   = "test"
)

type Result struct {
	id  int
	age int
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connect to db success!")
	}

	sql := "select * from test"
	rows, errs := db.Query(sql)

	if errs != nil {
		log.Fatal(errs)
	}

	var result Result

	for rows.Next() {
		err := rows.Scan(&result.id, &result.age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("The ID is: ")
		fmt.Println(result.id)
		fmt.Print("The age is: ")
		fmt.Println(result.age)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

}
