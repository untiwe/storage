package cache

import (
	"encoding/json"

	"github.com/dgraph-io/ristretto"

	"storage/config"
	"storage/conventions"

)

var (
	cache *ristretto.Cache //кеш для хранения заказов
	keys  map[string]bool  //список всех ключей кеша, т.к.  он не умеет предоставлять их сам
)

// Создаем кеш
func Init() {

	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,                          // количество счетчиков для отслеживания частоты использования элементов
		MaxCost:     config.GetInt64("max-cache"), // максимальный объем кеша 10Мб
		BufferItems: 128,                          // количество элементов в буфере для уменьшения блокировок
	})
	if err != nil {
		panic(err)
	}
	keys = make(map[string]bool)

}

func GetByID(id string) (order conventions.Order) {
	value, ok := cache.Get(id)
	if !ok {
		//помечаем ключи, которых уже нет в кеше
		return conventions.Order{OrderUID: ""}
	}

	return value.(conventions.Order)

}

// анмаршалим JSON и добавляем новый Order в кеш
func Add(str string) {

	var order conventions.Order
	if err := json.Unmarshal([]byte(str), &order); err != nil {
		println("Error add new order in cache: %s", err)
		return
	}

	ok := cache.Set(order.OrderUID, order, 1)
	if !ok {
		println("key is not added")
	}
	keys[order.OrderUID] = true
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

	//удаляем старые (не актуальные) ключи
	for _, oldKey := range oldKeys {
		delete(keys, oldKey)
	}

	return orders
}
