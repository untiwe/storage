package db

import (
	"encoding/json"
	"fmt"
	"log"
	"storage/conventions"

	_ "github.com/lib/pq"
)

func FillСache() {
	orders, err := FetchAllOrdersData()
	if err != nil {
		log.Fatal(err)
	}

	for _, order := range orders {
		value, err := json.Marshal(order)
		if err != nil {
			panic(err.Error())
		}
		kache.Add(string(value))
	}

}

func FetchAllOrdersData() ([]conventions.Order, error) {
	var orders []conventions.Order

	db, err := createConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Извлечение данных из таблицы orders
	rows, err := db.Query(`
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order conventions.Order
		err = rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Извлечение данных из таблиц deliveries, payments и items для каждого заказа
	for i := range orders {
		order := &orders[i]

		// Извлечение данных из таблицы deliveries
		err = db.QueryRow(`
			SELECT name, phone, zip, city, address, region, email
			FROM deliveries
			WHERE order_uid = $1
		`, order.OrderUID).Scan(
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
		)
		if err != nil {
			return nil, err
		}

		// Извлечение данных из таблицы payments
		err = db.QueryRow(`
			SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
			FROM payments
			WHERE order_uid = $1
		`, order.OrderUID).Scan(
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
		)
		if err != nil {
			return nil, err
		}

		// Извлечение данных из таблицы items
		itemRows, err := db.Query(`
			SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
			FROM items
			WHERE order_uid = $1
		`, order.OrderUID)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item conventions.Item
			err = itemRows.Scan(
				&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
			)
			if err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}

		err = itemRows.Err()
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

// Обработчик для получения данных всех заказов
func GetAllOrders() (data []byte, err error) {

	// Извлечение данных всех заказов из базы данных
	orders, err := FetchAllOrdersData()
	if err != nil {
		return nil, err
	}

	// Сериализация данных заказов в JSON
	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}

	return ordersJSON, nil

}

func InsertOrder(order conventions.Order) error {

	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer tx.Rollback() // Гарантируем откат транзакции в случае ошибки

	// Вставка данных в таблицу orders
	_, err = tx.Exec(`
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("Error insert data table orders: %v", err)
	}

	// Вставка данных в таблицу deliveries
	_, err = tx.Exec(`
		INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("Error insert data table deliveries: %v", err)
	}

	// Вставка данных в таблицу payments
	_, err = tx.Exec(`
		INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("Error insert data table payments: %v", err)
	}

	// Вставка данных в таблицу items
	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, item.ChrtID, order.OrderUID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("Error insert data table items: %v", err)
		}
	}

	// Фиксация транзакции
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error fix transaction: %v", err)
	}

	return nil
}
