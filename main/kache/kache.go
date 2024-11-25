package kache

import "sync"

// OrdersSet тип для хранения набора строк с ограничением строк
type OrdersSet struct {
	strings      []string
	maxSize      int
	currentIndex int
	mu           sync.Mutex
}

// NewOrdersSet создает новый набор строк с указанным максимальным размером
func NewOrdersSet(maxSize int) *OrdersSet {
	return &OrdersSet{
		strings:      make([]string, maxSize),
		maxSize:      maxSize,
		currentIndex: 0,
	}
}

// добавляем новый ordr в виде JSON
func (s *OrdersSet) Add(str string) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentIndex >= s.maxSize {
		s.currentIndex = 0
	}
	s.strings[s.currentIndex] = str
	// Иначе просто добавляем новую строку
	s.currentIndex++

}

// Возвращаем все Orders
func (s *OrdersSet) GetAll() []string {
	var filteredStrings []string
	for _, str := range s.strings {
		if str != "" { //фильтруем пустые строки
			filteredStrings = append(filteredStrings, str)
		}
	}
	return filteredStrings
}
