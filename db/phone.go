package db

import "database/sql"


func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err 
	}

	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (db *DB) Close() error {
	return db.db.Close()
}


func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil{
			return err
		}
	}
	
	return nil
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

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}

	err = createPhoneNumbersTable(db)

	if err != nil{
		return err
	}
		return db.Close()
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		ID SERIAL,
		value VARCHAR(255)
	)`

	_, err := db.Exec(statement)

	return err
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}

	resetDB(db, dbName)

	return db.Close()
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)

	if err != nil {
		return err
	}

	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATAVASE IF EXISTS" + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}