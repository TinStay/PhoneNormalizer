package main

import (
	"bytes"
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

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)

	defer db.Close()

	// Add phone numbers to table
	if err := db.Seed(); err != nil{
		panic(err)
	}

	phones, err := db.GetAllPhones()
	must(err)

	for _, p := range phones{
		fmt.Printf("Working on... %v\n", p.Number)
		// Format number
		number := normalize(p.Number)

		if number != p.Number {
			fmt.Println("Updating or removing...", number)
			// Find current number
			existing, err := db.FindPhone(number)
			must(err)

			if existing != nil {
				// delete number
				must(db.DeletePhone(p.ID))
			} else {
				// update number
				p.Number = number
				must(db.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}

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
