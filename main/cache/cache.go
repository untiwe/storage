package cache

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/ristretto"

	"storage/config"
	"storage/conventions"
)

var (
	cache *ristretto.Cache //кеш для хранения заказов
	keys  map[string]bool  //список всех ключей кеша, т.к.  он не умеет предоставлять их сам
)

// Создаем кеш
func init() {

	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,                              // количество счетчиков для отслеживания частоты использования элементов
		MaxCost:     int64(config.GetInt("max-cahe")), // максимальное количество элементов
		BufferItems: 128,                              // количество элементов в буфере для уменьшения блокировок
	})
	if err != nil {
		panic(err)
	}
	keys = make(map[string]bool)

}

// анмаршалим JSON и добавляем новый Order в кеш
func Add(str string) {

	var order conventions.Order
	if err := json.Unmarshal([]byte(str), &order); err != nil {
		println("Error add new order in cache: %s", err)
		return
	}
	ok := cache.Set(str, order, 1)
	if !ok {
		println("key is not added")
	}
	keys[str] = true
}

// Возвращаем все Orders
func GetAll() []conventions.Order {
	oldKeys := make([]string, 0)
	var orders = make([]conventions.Order, 0)

	for key, _ := range keys {
		value, ok := cache.Get(key)
		if !ok {
			//помечаем ключи, которых уже нет в кеше
			oldKeys = append(oldKeys, key)
		} else {
			orders = append(orders, value.(conventions.Order))
		}
	}

	fmt.Println(len(keys))
	fmt.Println(len(oldKeys))

	//удаляем старые (не актуальные) ключи
	for _, oldKey := range oldKeys {
		delete(keys, oldKey)
	}

	return orders
}
