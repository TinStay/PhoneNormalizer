package db

import "database/sql"


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