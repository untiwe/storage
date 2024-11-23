package db

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func WriteData(data string) {
	db, err := createConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqladd := `INSERT INTO teststrings(
		strcolum)
		VALUES ('` + data + `');`
	_, err = db.Exec(sqladd)
	if err != nil {
		log.Fatal(err)
	}
	kache.Add(data)
}

func FillСache() {
	db, err := createConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlget := `SELECT * FROM public.teststrings;`

	rows, err := db.Query(sqlget)
	if err != nil {
		log.Fatal(err)
	}

	var value string
	for rows.Next() {
		if err := rows.Scan(&value); err != nil {
			fmt.Errorf("Erros scan row: %v", err)
		}
		kache.Add(value)
	}

	// Проверка наличия ошибок после завершения чтения строк
	if err := rows.Err(); err != nil {
		fmt.Errorf("Error read rows: %v", err)
	}

}
