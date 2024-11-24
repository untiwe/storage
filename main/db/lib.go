package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func createConnectionString(targetUser bool, targetDB bool) (connStr string) {

	user := "postgres"
	pass := "postgres"
	dbname := ""

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "localhost"
	}

	// Если уже настроена база т.е. подключаемся к конкретной
	if targetUser {
		user = ownerName
		pass = ownerPass
	}
	if targetDB {
		dbname = " dbname=" + DbName
	}

	connStr = fmt.Sprintf("user=%s password=%s sslmode=disable host=%s port=5432%s", user, pass, dbURL, dbname)

	return connStr
}

func createConnection() (*sql.DB, error) {

	// Подключение к базе данных PostgreSQL
	connStr := createConnectionString(true, true)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
