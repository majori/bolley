package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	filename := os.Getenv("FILENAME") // TODO: Read filename somehow

	match, err := parseSpreadsheet(filename)
	if err == nil {
		saveMatch(db, match)
		fmt.Println(match)
	}
}
