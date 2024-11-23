package db

import (
	"log"

	_ "github.com/lib/pq"
)

func WriteData(data string) {
	db, err := createConnection()
	if err != nil {
		log.Fatal(err)
	}

	sqladd := `INSERT INTO teststrings(
		strcolum)
		VALUES ('` + data + `');`
	_, err = db.Exec(sqladd)
	if err != nil {
		log.Fatal(err)
	}
}
