package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Генерируем нужную строку подключения к базе
// targetUser ножно ли подключаться от пользователя
// targetDB нужно ли подключаться к базе проекта
func createConnectionString(targetUser bool, targetDB bool) (connStr string) {

	user := "postgres"
	pass := "postgres"
	dbname := ""

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "localhost"
	}

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

// Создаем подключение к нашей БД
func createConnection() (*sql.DB, error) {

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
