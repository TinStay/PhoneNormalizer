package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "phoneDB"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)

	// Open database
	db, err := sql.Open("postgres", psqlInfo)
	must(err)

	err = createDB(db, dbname)
	must(err)

	must(createPhoneNumbersTable(db))

	db.Close()
}

func createPhoneNumbersTable(db *sql.DB) error{
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		ID SERIAL,
		value VARCHAR(255)
	)`

	_, err := db.Exec(statement)

	return err
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)

	if err != nil {
		return err
	}

	return nil
}

func must(err error){
	if err != nil{
		panic(err)
	}
}

func normalize(phone string) string {
	var buf bytes.Buffer

	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

func normalizeRegex(phone string) string {
	regex := regexp.MustCompile("[\\D]")
	return regex.ReplaceAllString(phone, "")
}
