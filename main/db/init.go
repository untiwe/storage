package db

import (
	"database/sql"
	"fmt"
	"log"
	"storage/config"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type CacheInterface interface {
	Add(string)
	GetAll() []string
}

var kache CacheInterface
var DbName string

// создает базу данных, если она не существует
func createDatabaseIfNotExists(connStr string) error {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("create database " + DbName)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P04" { //Если БД уже есть - ок
			return nil
		}
		log.Fatal(err)
	}
	return nil
}

// createTables создает таблицы в PostgreSQL для модели данных
func createTables(db *sql.DB) error {
	createOrdersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		order_uid TEXT PRIMARY KEY,
		track_number TEXT,
		entry TEXT,
		locale TEXT,
		internal_signature TEXT,
		customer_id TEXT,
		delivery_service TEXT,
		shardkey TEXT,
		sm_id INTEGER,
		date_created TIMESTAMP,
		oof_shard TEXT
	);
	`

	createDeliveriesTable := `
	CREATE TABLE IF NOT EXISTS deliveries (
		order_uid TEXT PRIMARY KEY,
		name TEXT,
		phone TEXT,
		zip TEXT,
		city TEXT,
		address TEXT,
		region TEXT,
		email TEXT,
		FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
	);
	`

	createPaymentsTable := `
	CREATE TABLE IF NOT EXISTS payments (
		order_uid TEXT PRIMARY KEY,
		transaction TEXT,
		request_id TEXT,
		currency TEXT,
		provider TEXT,
		amount INTEGER,
		payment_dt BIGINT,
		bank TEXT,
		delivery_cost INTEGER,
		goods_total INTEGER,
		custom_fee INTEGER,
		FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
	);
	`

	createItemsTable := `
	CREATE TABLE IF NOT EXISTS items (
		chrt_id INTEGER,
		order_uid TEXT,
		track_number TEXT,
		price INTEGER,
		rid TEXT,
		name TEXT,
		sale INTEGER,
		size TEXT,
		total_price INTEGER,
		nm_id INTEGER,
		brand TEXT,
		status INTEGER,
		FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
	);
	`

	_, err := db.Exec(createOrdersTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(createDeliveriesTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(createPaymentsTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(createItemsTable)
	if err != nil {
		return err
	}

	return nil
}

func Init(k CacheInterface) {

	kache = k
	DbName = config.GetString("db-name")

	connStr := createConnectionString(false)
	//создаем БД (если нету)
	createDatabaseIfNotExists(connStr)

	//добавляем пожклчюение БД
	connStr += " dbname=" + config.GetString("db-name")
	// Подключение к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблиц
	err = createTables(db)
	if err != nil {
		log.Fatal(err)
	}
}
