package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"regexp"

	// phonedb "github.com/TinStay/PhoneNormalizer/db"
	phonedb "github.com/Basics/src/github.com/TinStay/PhoneNormalizer/db"

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
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	// Migrate database
	must(phonedb.Migrate("postgres", psqlInfo))

	db, err := sql.Open("postgres", psqlInfo)
	must(err)

	defer db.Close()

	// Add phone to table
	_, err = insertPhone(db, "1234567890")
	must(err)
	_, err = insertPhone(db, "123 456 7891")
	must(err)
	_, err = insertPhone(db, "(123) 456 7892")
	must(err)
	_, err = insertPhone(db, "(123) 456-7893")
	must(err)
	_, err = insertPhone(db, "123-456-7890")
	must(err)
	_, err = insertPhone(db, "1234567892")
	must(err)
	_, err = insertPhone(db, "(123)456-7892")
	must(err)

	// number, err := getPhone(db, id)
	// must(err)
	// fmt.Println("Number is: ", number)

	phones, err := getAllNumbers(db)
	must(err)

	for _, p := range phones{
		fmt.Printf("Working on... %v\n", p.number)
		// Format number
		number := normalize(p.number)

		if number != p.number {
			fmt.Println("Updating or removing...", number)
			// Find current number
			existing, err := findPhone(db, number)
			must(err)

			if existing != nil {
				// delete number
				must(deletePhone(db, p.id))
			} else {
				// update number
				p.number = number
				must(updatePhone(db, p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}

	// Close db
	db.Close()
}

type phone struct{
	id int
	number string
}

func updatePhone(db *sql.DB, p phone) error {
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.Exec(statement, p.id, p.number)
	return err
}

func deletePhone(db *sql.DB, id int) error {
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.Exec(statement, id)
	return err
}

func findPhone(db *sql.DB, number string)(*phone, error){
	var p phone

	row := db.QueryRow("Select * FROM phone_numbers WHERE value=$1", number)
	err := row.Scan(&p.id, &p.number) 

	if err != nil {
		if err == sql.ErrNoRows{
			return nil, nil
		}else {
			return nil, err
		}
	}

	return &p, nil
}

func getAllNumbers(db *sql.DB)([]phone, error){
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ret []phone

	// Loop through all rows
	for rows.Next(){
		var p phone
		if err := rows.Scan(&p.id,&p.number,); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func getPhone(db *sql.DB, id int)(string, error){
	var number string

	row := db.QueryRow("Select value FROM phone_numbers WHERE id=$1", id)
	err := row.Scan(&number) 

	if err != nil {
		return "", err
	}

	return number, nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`

	var id int
	// Add phone number to table
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil{
		return -1, err
	}

	return int(id), nil
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
