package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const DbName = "shared"

func createConnectionString(target bool) (connStr string) {

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "localhost"
	}
	// Параметры подключения к базе данных
	connStr = "user=postgres sslmode=disable password=postgres host=" + dbURL + " port=5432"

	if target {
		dbParam := " dbname=" + DbName
		connStr += dbParam
	}
	return connStr
}

func createConnection() (*sql.DB, error) {

	// Подключение к базе данных PostgreSQL
	connStr := createConnectionString(true)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка соединения с базой данных: %v", err)
	}

	return db, nil
}
